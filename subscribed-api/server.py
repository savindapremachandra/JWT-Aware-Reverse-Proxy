from flask import Flask
app = Flask(__name__)

@app.route("/")
def subscribed_home():
    return {"msg": "Hello Subscribed User 🚀"}

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=9000)