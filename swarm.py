import asyncio
import base64
import io
import time
from uuid import uuid4

import aiohttp
from faker import Faker
from PIL import Image

fake = Faker()


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


async def device_swarm(num_devices=100):
    async with aiohttp.ClientSession() as session:
        tasks = []
        for i in range(1, num_devices + 1):
            task = asyncio.ensure_future(
                send_post_request(session, device_id=str(uuid4()))
            )
            tasks.append(task)
        await asyncio.gather(*tasks)


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    for i in range(100):
        loop.run_until_complete(device_swarm())
        time.sleep(1)
