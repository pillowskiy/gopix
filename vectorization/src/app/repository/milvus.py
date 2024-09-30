import os
import numpy as np
from pymilvus import connections, FieldSchema, CollectionSchema, DataType, Collection, utility

MILVUS_HOST = os.getenv("MILVUS_HOST", "127.0.0.1")
MILVUS_PORT = os.getenv("MILVUS_PORT", "19530")

connections.connect(host=MILVUS_HOST, port=MILVUS_PORT)

def create_collection_if_not_exists(collection_name: str, dim: int) -> Collection:
   if utility.has_collection(collection_name):
        return Collection(collection_name)

   return _create_clip_collection(collection_name, dim)

def create_collection_force(collection_name: str, dim: int) -> Collection:
    if utility.has_collection(collection_name):
        utility.drop_collection(collection_name)

    return _create_clip_collection(collection_name, dim)

def _create_clip_collection(collection_name: str, dim: int) -> Collection:
    print("[Milvus]: Creating clip collection")

    schema = CollectionSchema(fields=[
        FieldSchema(name='id', dtype=DataType.INT64, is_primary=True, auto_id=False),
        FieldSchema(name='vector', dtype=DataType.FLOAT_VECTOR, dim=dim)
    ])

    print(f"[Milvus]: Define colection with name {collection_name}")
    collection = Collection(name=collection_name, schema=schema)

    print(f"[Milvus]: Create indexes for collection {collection_name}")
    collection.create_index(field_name="vector", index_params={
        'metric_type': 'L2',
        'index_type': "IVF_FLAT",
        'params': {"nlist": dim}
    })

    return collection

class MilvusRepository:
    def __init__(self, collection_name: str, dim: int, force_create: bool = False):
        create_collection = create_collection_force if force_create is True else create_collection_if_not_exists
        create_collection(collection_name, dim)

        self.collection = Collection(collection_name)
        self.collection.load()

    def insert(self, id: int, vector: np.ndarray):
        data = [{"id": id, "vector": vector}]
        self.collection.insert(data=data)

    def search_neighbors(self, vector: np.ndarray, limit: int = 20):
        results = self.collection.search(
            anns_field="vector",
            data=[vector],
            param={
                "metric_type": "L2",
                "params": {
                    "nprobe": 10
                },
            },
            limit=limit,
            output_fields=["id"],
        )

        return results[0]

    def delete(self, id_value: int):
        self.collection.delete(expr=f"id == {id_value}")

    def get_vector_by_id(self, id_value: int):
        query_result = self.collection.query(
            expr=f"id == {id_value}",
            output_fields=["vector"],
        )

        if not query_result:
            raise ValueError(f"No vector found with ID {id_value}")

        vector = query_result[0]["vector"]
        return vector

    def exists(self, id_value: int):
        query_result = self.collection.query(
            expr=f"id == {id_value}",
            output_fields=["id"],
        )
        return bool(query_result)
