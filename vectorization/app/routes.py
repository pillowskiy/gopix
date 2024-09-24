from flask import Blueprint, jsonify, request

from .milvus import MilvusClient
from .model import extract_features

main = Blueprint("main", __name__)
milvus_client = MilvusClient(collection_name="l2")


@main.route("/features", methods=["POST"])
def extract_features_endpoint():
    try:
        file = request.files["image"]
        target_id = int(request.form["id"])
        vector = extract_features(file)
        milvus_client.insert_vector(target_id=target_id, vector=vector)

        return jsonify({"vector": vector.tolist(), "id": target_id})
    except Exception as e:
        raise e
        return jsonify({"error": str(e)}), 500


@main.route("/similar/<int:id>", methods=["GET"])
def search_neighbors_endpoint(id):
    try:
        limit = int(request.args.get("limit", 10))
        neighbors = milvus_client.search_neighbors_by_id(id_value=id, limit=limit)
        results = [
            {"id": result["id"], "distance": result["distance"]} for result in neighbors
        ]
        return jsonify(results), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 400
    except Exception as e:
        return jsonify({"error": str(e)}), 500
