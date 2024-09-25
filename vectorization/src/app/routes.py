from flask import Blueprint, jsonify, request

from .milvus import MilvusClient
from .model import extract_features

main = Blueprint("main", __name__)
milvus_client = MilvusClient(collection_name="l2")


@main.route("/features", methods=["POST"])
def extract_features_endpoint():
    try:
        id_str = request.form.get("id")
        if id_str is None:
            return jsonify({"error": "ID not provided"}), 400

        try:
            target_id = int(id_str)
        except ValueError:
            return jsonify({"error": "Invalid ID format"}), 400

        if "image" not in request.files:
            return jsonify({"error": "No file provided"}), 400

        if milvus_client.exists(target_id):
            return jsonify({"error": "Vector with the same ID already exists"}), 409

        file = request.files["image"]
        vector = extract_features(file)
        milvus_client.insert_vector(target_id=target_id, vector=vector)

        return jsonify({"vector": vector.tolist(), "id": target_id})
    except Exception as e:
        raise e
        return jsonify({"error": str(e)}), 500


@main.route("/similar/<int:id>", methods=["GET"])
def search_neighbors_endpoint(id):
    try:
        limit = int(request.args.get("limit", 20))
        vector = milvus_client.get_by_id(id)
        neighbors = milvus_client.search_neighbors(vector, limit)

        results = [
            {"id": result["id"], "distance": result["distance"]} for result in neighbors
        ]
        return jsonify(results), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 400
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@main.route("/features/<int:id>", methods=["DELETE"])
def delete_by_id_endpoint(id):
    try:
        if not milvus_client.exists(id):
            return jsonify({"error": "Vector not found"}), 404

        milvus_client.delete(id)
        return jsonify({"success": True}), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 400
    except Exception as e:
        return jsonify({"error": str(e)}), 500
