import asyncio
import base64
import io
import random
import time
from uuid import uuid4

import aiohttp
from PIL import Image

DEVICE_COUNT = 100


def dummy_bytes():
    image = Image.new("RGB", (100, 100))
    image_bytes = io.BytesIO()
    image.save(image_bytes, format="JPEG")
    return base64.b64encode(image_bytes.getvalue()).decode()


async def send_post_request(session, device_id):
    url = "http://127.0.0.1:5555/buff"
    data = {"deviceId": device_id, "data": dummy_bytes()}
    async with session.post(url, json=data) as response:
        print(
            f"Device {device_id}: Status Code {response.status} Response: {await response.text()}"
        )


async def device_swarm():
    num_devices = random.randint(0, DEVICE_COUNT)
    print(f"Device cnt {num_devices}")
    async with aiohttp.ClientSession() as session:
        tasks = []
        for i in range(num_devices):
            task = asyncio.ensure_future(
                send_post_request(session, device_id=str(uuid4()))
            )
            tasks.append(task)
        await asyncio.gather(*tasks)


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    while True:
        loop.run_until_complete(device_swarm())
        time.sleep(5)
