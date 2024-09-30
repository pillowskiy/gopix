import os

import clip
import numpy as np
import torch
from PIL import ImageFile

CLIP_DOWNLOAD_ROOT = os.getenv("CLIP_DOWNLOAD_ROOT", "./models")

device = "cuda" if torch.cuda.is_available() else "cpu"
print(
    f"[ClipModel]: CUDA {'is' if torch.cuda.is_available() else "isn't"} available. Default device: {device.upper()}"
)

model, preprocess = clip.load(
    "ViT-B/32",
    device=device,
    download_root=CLIP_DOWNLOAD_ROOT,
)


class ClipModel:
    def vectorize_text(self, text: str, device: str = device) -> np.ndarray:
        text_input = clip.tokenize([text]).to(device)
        with torch.no_grad():
            vec = model.encode_text(text_input)
            vec = vec / vec.norm(dim=-1, keepdim=True)
        return vec.cpu().numpy().squeeze()

    def vectorize_image(self, img: ImageFile, device: str = device) -> np.ndarray:
        img = preprocess(img).unsqueeze(0).to(device)
        with torch.no_grad():
            vec = model.encode_image(img)
            vec = vec / vec.norm(dim=-1, keepdim=True)

        return vec.cpu().numpy().squeeze()
