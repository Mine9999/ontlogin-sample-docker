<script setup>
import {
  createAuthRequest,
  postRequest,
  requestQR,
  queryQRResult,
} from "ontlogin";
import QRCode from "qrcode-svg";

let svg = "";
const showQr = (text) => {

  console.log("show qr code of", text);
  svg = new QRCode({
        content: text,
        width: 122,
        height: 122,
      }).svg();
  
  document.getElementById("svg_container").innerHTML = svg;
};

const login = async () => {
  const authRequest = createAuthRequest(0); 
  console.log("authRequest", authRequest);
  const authChallenge = await postRequest("http://localhost:3000/requestChallenge", authRequest );
  console.log("authChallenge", authChallenge);
  const { text, id } = await requestQR(authChallenge);
  showQr(text);
  const challengeResponse = await queryQRResult(id);
  console.log("challengeResponse", challengeResponse);
  const result = await fetch("http://localhost:3000/submitChallenge", {
            method: "post",
            body: JSON.stringify(challengeResponse),
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
        }).then((res) => {
            return res;
        }).catch(err => console.log(err));
        console.log("result", result);

        let crDid = challengeResponse.did;
        if(result.ok){
          document.getElementById("svg_container").innerHTML = "";
          document.getElementById("submit_did").innerHTML = crDid;
        }
};
</script>

<template>
  <div>
    <button @click="login">sign in with ONT LOGIN</button>
    <div id="svg_container"/>
    <div id="submit_did"/>
  </div>
</template>
