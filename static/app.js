
function init()
{
  document.getElementById("registerBtn").onclick = register;
  document.getElementById("loginBtn").onclick = login;
  console.log("done setup onclick event");
}

function register() {
  console.log("register!")
  let msg = document.getElementById("msg");
  msg.innerHTML = "<div class='lds-dual-ring'></div>";
  msg.classList.remove("hidden");

  startRegister().then(function() {
    console.log("done register!")
    return false;
  }).catch(function() {
    console.log("failed to startRegister");
  });
}

async function startRegister()
{
  console.log("startRegister!")
  let resp = await fetch("/register")
  let u2fChallenge = await resp.json()
  console.log("register result:")
  console.log(u2fChallenge)
  u2f.register(u2fChallenge.appId,
               u2fChallenge.registerRequests,
               u2fChallenge.registeredKeys,
               u2fRegisterCallback,
               30);
  console.log("StartRegister end")
}

async function u2fRegisterCallback(challengeSigData)
{
  console.log("u2fRegister callback!")
  console.log("param:")
  console.log(challengeSigData)
  console.log("as string: "+JSON.stringify(challengeSigData));
  let resp = await fetch("/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(challengeSigData)
  })
  console.log("The response from the challengeSignature Request")
  console.log(resp)
  console.log("as json:")
  let j = await resp.json()
  console.log(j)
  document.getElementById("msg").innerHTML = "registered: "+JSON.stringify(j);
  console.log("registercallback finish")
}

function login() {
  console.log("Login!");
  let msg = document.getElementById("msg");
  msg.innerHTML = "<div class='lds-dual-ring'></div>";
  msg.classList.remove("hidden");
  startLogin().then(function() {
    console.log("done login!")
    return false;
  }).catch(function() {
  	console.log("failed to startRegister");
  });
  return false;
}

async function startLogin()
{
  console.log("startLogin!")
  let resp = await fetch("/login")
  let signChallenge = await resp.json()
  console.log("login result:")
  console.log(signChallenge)
  u2f.sign(signChallenge.appId,
    signChallenge.challenge,
    signChallenge.registeredKeys,
    u2fLoginCallback,
    30)
  console.log("StartLogin end")
}
async function u2fLoginCallback(loginChallengeData) {
  console.log("u2fLoginCallback")
  console.log(loginChallengeData)
  console.log("as string: "+JSON.stringify(loginChallengeData));
  let resp = await fetch("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(loginChallengeData)
  })
  console.log("The response from the loginChallengeSign Request")
  console.log(resp)
  console.log("as json:")
  let j = await resp.json()
  console.log(j)
  document.getElementById("msg").innerHTML = "login: "+JSON.stringify(j);
  console.log("loginCallback finish")

}
