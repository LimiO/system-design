import json

from flask import Flask

app = Flask(__name__)

@app.route("/")
def index():
    return "<p>Hello, World!</p>"


@app.route("/health/")
def health():
    return json.dumps({'status':"ok"}), 200, {'ContentType':'application/json'} 



if __name__ == '__main__':
   app.run(host='0.0.0.0', port=8000)