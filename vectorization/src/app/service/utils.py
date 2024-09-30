from flask import Response, jsonify
from typing import Dict, Any

class ServiceError:
    _message: str
    _status: int

    def __init__(self, message: str, status: int):
        self._message = message
        self._status = status

    def to_flask_res(self) -> tuple[Response, int]:
        return jsonify(self.to_dict()), self._status

    def to_dict(self) -> Dict[str, Any]:
        return { "message": self._message }