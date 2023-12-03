import requests

url = "http://localhost:8085"


def test_flow():
    headers = {}
    username = "mytestuser3"
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
        "amount": 100,
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
    assert r.json()["balance"] == 40

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
    assert r.json()["balance"] == 40


if __name__ == "__main__":
    test_flow()
