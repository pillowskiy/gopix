import time
import os
from flask import Blueprint, jsonify, request

from .repository.milvus import MilvusRepository
from .model.clip import ClipModel
from .service.utils import ServiceError
from .service.vectorization import VectorizationService

main = Blueprint("main", __name__)

model = ClipModel()

milvus_collection = os.getenv("MILVUS_COLLECTION", "clip_collection")
milvus_dim = os.getenv("MILVUS_DIM", 512)
repo = MilvusRepository(milvus_collection, milvus_dim)

service = VectorizationService(model, repo)

@main.route("/features", methods=["POST"])
def extract_features_endpoint():
    try:
        try:
            id_str = request.form.get("id")
            if id_str is None: raise ValueError("invalid")
            target_id = int(id_str)
        except ValueError:
            return jsonify({"error": "Invalid ID format"}), 400

        if "image" not in request.files:
            return jsonify({"error": "No file provided"}), 400

        file = request.files["image"]

        start = time.perf_counter()
        service.insert(file, target_id)
        took = f"Featurize and insert took: {(time.perf_counter() - start)*1000:2f}ms"

        return jsonify({ "message": took }), 201
    except ServiceError as e:
        return e.to_flask_res()
    except Exception as e:
        return jsonify({"message": str(e)}), 500



@main.route("/similar/<int:id>", methods=["GET"])
def get_similar_endpoint(id):
    try:
        limit = int(request.args.get("limit", 20))
        results = service.search_similar(id, limit)

        return jsonify(results), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 400
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@main.route("/search", methods=["GET"])
def search_by_text_endpoint():
    try:
        limit = int(request.args.get("limit", 20))
        text = request.args.get("query")
        results = service.search_by_text(text, limit)

        return jsonify(results), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 400
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@main.route("/features/<int:id>", methods=["DELETE"])
def delete_by_id_endpoint(id):
    try:
        service.delete(id)
        return jsonify({"success": True}), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 400
    except ServiceError as e:
        return e.to_flask_res()
    except Exception as e:
        return jsonify({"error": str(e)}), 500
