const { readFileSync } = require("fs");
const privateKey = readFileSync("./private.json", "utf8");

let jose = require("node-jose");

let header = {
  alg: "RS256",
  typ: "JWT",
  kid: process.env.KID,
};

console.log(process.env.KID);

let payload = {
  iss: "2004337545",
  sub: "2004337545",
  aud: "https://api.line.me/",
  exp: Math.floor(new Date().getTime() / 1000) + 60 * 30,
  token_exp: 60 * 60 * 24 * 30,
};

jose.JWS.createSign(
  { format: "compact", fields: header },
  JSON.parse(privateKey)
)
  .update(JSON.stringify(payload))
  .final()
  .then((result) => {
    console.log(result);
  });
