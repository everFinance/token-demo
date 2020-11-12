let provider
let signer
let account
const rollURL = "../"

const lessAddress = (addr) => {
    return addr.substr(0, 6) + "..." + addr.substr(38, 4)
}

const signTx = async (tx) => {
    const signData = tx.nonce + tx.type + tx.from + tx.to + tx.amount
    return signer.signMessage(signData)
}

const transferToken = async () => {
    const url = rollURL + "tx"

    const tx = {
        nonce: Date.parse(new Date().toString()).toString(),
        type: "transfer",
        from: account,
        to: $("#target").val(),
        amount: $("#target_amount").val(),
        sign: ""
    }
    tx.sign = await signTx(tx)

    const resp = await fetch(url, {
        body: JSON.stringify(tx),
        headers: {
            'content-type': 'application/json'
        },
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
    })

    const jResp = await resp.json()
    if (jResp.error) {
        alert(jResp.error)
    } else {
        alert("转账成功")
    }
}

const claimToken = async () => {
    const url = rollURL + "tx"

    const tx = {
        nonce: Date.parse(new Date().toString()).toString(),
        type: "claim",
        from: account,
        to: "",
        amount: "100",
        sign: ""
    }
    tx.sign = await signTx(tx)

    const resp = await fetch(url, {
        body: JSON.stringify(tx),
        headers: {
            'content-type': 'application/json'
        },
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
    })

    const jResp = await resp.json()
    if (jResp.error) {
        alert(jResp.error)
    } else {
        alert("获得 100 ROL")
    }
}

const genTx = (tx) => {
    let txType // 转入，转出，领取
    let status = {}
    let message
    let amount

    if (tx.type == "claim") {
        txType = "Claim"
        message = "New user claimed"
        amount = tx.amount
    } else if(tx.type = "transfer") {
        amount = tx.amount
        message = lessAddress(tx.from) + " to " + lessAddress(tx.to)
        if(tx.from == account) {
            txType = "Out"
        } else if(tx.to == account){
            txType = "In"
        }
    } else {
        txType = "invalid"
    }

    if (tx.id != "") {
        status.message = "Packaged"
        status.color = "#6f42c1" // 紫色
        status.url = "https://viewblock.io/arweave/tx/" + tx.id
        status.target = `target="_blank"`
    } else {
        status.message = "Confirmed"
        status.color = "#007bff" // 蓝色
        status.url = "#"
        status.target = ""
    }

    return `
    <div class="media text-muted pt-3">
        <svg class="bd-placeholder-img mr-2 rounded" width="32" height="32"
            xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid slice" focusable="false"
            role="img" aria-label="Placeholder: 32x32">
            <title>Placeholder</title>
            <rect width="100%" height="100%" fill="` + status.color + `" />
        </svg>
        <div class="media-body pb-3 mb-0 small lh-125 border-bottom border-gray">
            <div class="d-flex justify-content-between align-items-center w-100">
                <strong class="text-gray-dark">` + txType + `</strong>
                <a href="` + status.url + `" ` + status.target + `>` + status.message + `</a>
            </div>
            <span class="d-block">` + message + `</span>
            <span>Amount: ` + amount + `</span>
        </div>
    </div>
    `
}

const genTxs = (txs) => {
    let html = ""

    for(tx of txs) {
        html += genTx(tx)
    }

    return html
}

const updateBalance = async () => {
    const url = rollURL + "balanceOf/" + account

    const resp = await fetch(url, { method: 'GET' })
    const jResp = await resp.json()

    $("#amount").html(jResp.balance)
}

const updateTxs = async () => {
    const url = "../txs/" + account

    const resp = await fetch(url, { method: 'GET' })
    const jResp = await resp.json()
    let txs = jResp.txs
    txs.reverse()

    $("#txs").html(genTxs(txs))
}

const updateInfo = async() => {
    $("#connect").html(lessAddress(account))

    updateBalance()
    updateTxs()
}

const connectWallet = async () => {
    // await window.ethereum.enable()
    await ethereum.request({ method: 'eth_requestAccounts' })
    provider = new ethers.providers.Web3Provider(ethereum)
    signer = provider.getSigner()
    account = await signer.getAddress()

    updateInfo()
}

$("#connect").click(async (e) => {
    e.preventDefault()
    if (typeof window.ethereum === 'undefined') {
        alert("Wallet is not installed!");
        return
    }

    connectWallet()
})

$("#transfer").click(async (e) => {
    e.preventDefault()

    await transferToken()
    updateInfo()
})


$("#claim").click(async (e) => {
    e.preventDefault()

    await claimToken()
    updateInfo()
})

$('#reflesh_txs').click(async (e) => {
    e.preventDefault()

    updateInfo()
})

ethereum.on('accountsChanged', async (a) => {
    connectWallet()
});
