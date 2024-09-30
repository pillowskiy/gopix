import io
import os
from PIL import Image
from ..model.clip import ClipModel
from ..repository.milvus import MilvusRepository
from .utils import ServiceError

from typing import IO, Union

ImageFile = Union[str, bytes, "os.PathLike[str]", "os.PathLike[bytes]"] | IO[bytes]

class VectorizationService:
    _model: ClipModel
    _repo: MilvusRepository

    def __init__(self, model: ClipModel, repo: MilvusRepository):
        self._model = model
        self._repo = repo

    def insert(self, image_file: ImageFile, image_id: str):
        if self._repo.exists(image_id):
            raise ServiceError(f"Vector with key {image_id} already exists", 409)

        img = Image.open(io.BytesIO(image_file.read()))
        vec = self._model.vectorize_image(img)
        self._repo.insert(image_id, vec)


    def search_by_text(self, text: str, limit: int = 5):
        vec = self._model.vectorize_text(text)
        neighbors = self._repo.search_neighbors(vec, limit)
        return [result.id for result in neighbors]

    def search_similar(self, image_id: int, limit: int = 5):
        vec = self._repo.get_vector_by_id(image_id)
        neighbors = self._repo.search_neighbors(vec, limit)
        return [result.id for result in neighbors]

    def delete(self, image_id: int):
        if not self._repo.exists(image_id):
            raise ServiceError("Vector not found", 404)

        self._repo.delete(image_id)
