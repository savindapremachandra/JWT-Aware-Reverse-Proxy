from flask import Flask
app = Flask(__name__)

@app.route("/")
def free_home():
    return {"msg": "Hello Free-tier User 🌱"}

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=9000)