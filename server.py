import sys

from flask import Flask

app = Flask(__name__)


@app.route("/")
def hello_world():
    return f"<p>Hello from {sys.argv[1]}!</p>"


if __name__ == "__main__":
    app.run(port=sys.argv[2])
