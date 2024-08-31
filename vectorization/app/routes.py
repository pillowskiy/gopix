from flask import Blueprint, jsonify, request

from .model import extract_features

main = Blueprint("main", __name__)


@main.route("/features", methods=["POST"])
def extract_features_endpoint():
    try:
        file = request.files["image"]
        vector = extract_features(file)
        return jsonify({"vector": vector.tolist()})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
