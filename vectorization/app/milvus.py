import os

import numpy as np
from pymilvus import CollectionSchema, DataType, FieldSchema
from pymilvus import MilvusClient as MC

MILVUS_HOST = os.getenv("MILVUS_HOST", "localhost")
MILVUS_PORT = os.getenv("MILVUS_PORT", "19530")
MILVUS_FILE = os.getenv("MILVUS_FILE", "./milvus_data/milvus.db")


class MilvusClient:
    def __init__(self, collection_name: str):
        self.collection_name = collection_name
        self.client = MC(MILVUS_FILE)
        self._load_collection()

    def _load_collection(self):
        if self.client.has_collection(self.collection_name):
            return

        fields = [
            FieldSchema(
                name="id", dtype=DataType.INT64, is_primary=True, auto_id=False
            ),
            FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=512),
        ]
        schema = CollectionSchema(fields=fields)

        self.client.create_collection(
            collection_name=self.collection_name,
            schema=schema,
            metric_type="L2",
        )

        index_params = MC.prepare_index_params()
        index_params.add_index(
            field_name="vector",
            index_name="vector_index",
            metric_type="L2",
            params={"nlist": 16384},
        )

        self.client.create_index(
            collection_name=self.collection_name,
            index_params=index_params,
        )

    def insert_vector(self, target_id: int, vector: np.ndarray):
        if vector.shape[0] != 512:
            raise ValueError("Vector dimension must be 512.")

        data = [{"id": target_id, "vector": vector}]
        self.client.insert(
            collection_name=self.collection_name,
            data=data,
        )
        print(f"Inserted vector with ID {target_id}.")

    def search_neighbors_by_id(self, id_value: int, limit: int = 10):
        query_result = self.client.query(
            collection_name=self.collection_name,
            filter=f"id == {id_value}",
            output_fields=["vector"],
        )

        if not query_result:
            raise ValueError(f"No vector found with ID {id_value}")

        vector = query_result[0]["vector"]

        results = self.client.search(
            collection_name=self.collection_name,
            anns_field="vector",
            data=[vector],
            params={"nprobe": 128},
            metric_type="L2",
            limit=limit,
            output_fields=["id"],
        )

        return results[0]
