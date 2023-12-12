import requests
from random import choice
from string import ascii_letters

url = "http://localhost:8085"


def get_id():
    return "".join(choice(ascii_letters) for _ in range(20))


def test_flow():
    headers = {}
    username = get_id()
    data = {
        "username": username,
        "first_name": "test",
        "last_name": "user",
        "email": "testuser@gmail.com",
        "phone": 432479238472389,
        "password": "mybestpassword",
    }
    r = requests.post(url + "/user", headers=headers, json=data)
    assert r.status_code == 200

    token_data = r.json()
    headers = {
        "Authorization": "Bearer {}".format(token_data["token"])
    }
    data = {
        "username": username,
        "amount": 400,
    }
    r = requests.post(url + "/billing/balance/add", headers=headers, json=data)
    assert r.status_code == 200

    data = {
        "count": 2,
        "price": 30,
        "product_id": 1,
    }
    r = requests.post(url + "/buy", headers=headers, json=data)
    assert r.status_code == 200
    order_id = r.json()["order_id"]
    total = r.json()["total_price"]
    assert total == 60

    data = {
        "username": username,
    }
    r = requests.get(url + "/billing/balance", headers=headers, json=data)
    assert r.status_code == 200
    assert r.json()["balance"] == 340

    data = {
        "count": 2,
        "price": 30,
        "product_id": 1,
    }
    r = requests.post(url + "/buy", headers=headers, json=data)
    assert r.status_code == 200

    data = {
        "username": username,
    }
    r = requests.get(url + "/billing/balance", headers=headers, json=data)
    assert r.status_code == 200
    assert r.json()["balance"] == 280

    # zero cours
    data = {
        "count": 2,
        "price": 30,
        "product_id": 1,
    }
    r = requests.post(url + "/buy", headers=headers, json=data)
    assert r.status_code != 200

    data = {
        "count": 50,
    }
    r = requests.get(url + "/orders", headers=headers, json=data)
    assert r.status_code == 200
    last_order = r.json()["orders"][-1]
    assert last_order["paid"] == 2

    data = {
        "username": username,
    }
    r = requests.get(url + "/billing/balance", headers=headers, json=data)
    assert r.status_code == 200
    assert r.json()["balance"] == 280



def test_stock():
    headers = {}
    username = get_id()
    data = {
        "username": username,
        "first_name": "test",
        "last_name": "user",
        "email": "testuser@gmail.com",
        "phone": 432479238472389,
        "password": "mybestpassword",
    }
    r = requests.post(url + "/user", headers=headers, json=data)
    assert r.status_code == 200

    token_data = r.json()
    headers = {
        "Authorization": "Bearer {}".format(token_data["token"])
    }
    data = {
        "product_id": 0,
    }
    r = requests.get("http://localhost:8086/product", headers=headers, json=data)

    assert r.status_code == 200
    count = r.json()["count"]

    data = {
        "product_id": 0,
        "count": 2
    }
    r = requests.post("http://localhost:8086/products/reserve", headers=headers, json=data)
    assert r.status_code == 200
    reserve_id = r.json()["reserve_id"]

    data = {
        "product_id": 0,
    }
    r = requests.get("http://localhost:8086/product", headers=headers, json=data)

    assert r.status_code == 200
    assert r.json()["count"] == count - 2

    data = {
        "reserve_id": reserve_id,
        "status": 1,
    }
    r = requests.post("http://localhost:8086/products/commit", headers=headers, json=data)
    assert r.status_code == 200


def test_courier():
    headers = {}
    username = get_id()
    data = {
        "username": username,
        "first_name": "test",
        "last_name": "user",
        "email": "testuser@gmail.com",
        "phone": 432479238472389,
        "password": "mybestpassword",
    }
    r = requests.post(url + "/user", headers=headers, json=data)
    assert r.status_code == 200

    token_data = r.json()
    headers = {
        "Authorization": "Bearer {}".format(token_data["token"])
    }
    r = requests.post("http://localhost:8087/courier/reserve", headers=headers)
    assert r.status_code == 200
    courusername = r.json()['username']

    data = {
        "username": courusername,
    }
    r = requests.post("http://localhost:8087/courier/unreserve", headers=headers, json=data)
    assert r.status_code == 200


def test_stock_and_cour():
    headers = {}
    username = get_id()
    data = {
        "username": username,
        "first_name": "test",
        "last_name": "user",
        "email": "testuser@gmail.com",
        "phone": 432479238472389,
        "password": "mybestpassword",
    }
    r = requests.post(url + "/user", headers=headers, json=data)
    assert r.status_code == 200

    token_data = r.json()
    headers = {
        "Authorization": "Bearer {}".format(token_data["token"])
    }


# TODO(albert-si) на завтра. Проверить весь флоу в backend просто
if __name__ == "__main__":
    pass
    test_flow()
    # test_courier()
    # test_stock()
    # test_stock()