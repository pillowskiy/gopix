import io
import os

import torch
import torchvision
from PIL import Image
from torchvision import transforms

os.environ["TORCH_HOME"] = os.path.join(os.getcwd(), "models")

model = torchvision.models.resnet18(weights="DEFAULT")
model.eval()

transform = transforms.Compose(
    [
        transforms.Resize((256, 256)),
        transforms.ToTensor(),
        transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
    ]
)


def extract_features(image_file):
    img = Image.open(io.BytesIO(image_file.read())).convert("RGB")
    img = transform(img)

    with torch.no_grad():
        activation = {}

        def hook(model, input, output):
            activation["avgpool"] = output.detach()

        model.avgpool.register_forward_hook(hook)
        _ = model(img[None, ...])

        vec = activation["avgpool"].numpy().squeeze()

    return vec
