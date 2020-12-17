const ecc = require('eosjs-ecc')
/*
ecc.randomKey().then(privateKey => {
    console.log('Private Key:\t', privateKey) // wif
    console.log('Public Key:\t', ecc.privateToPublic(privateKey)) // EOSkey...
})
*/

// SIG_K1_Jz4KTq5v3dhcRKYehza6NRF5SxZaEzdPiVohRhX5SqoDCjmjf6hh3vyqfFHzEUagWZHC3L6G6SaJvKdH3UWEVwJRLWc6jL
    let f = async ()=>{
        let priv = '5K8bmn8AMNewSzgB3VNnz7pahVVLTF7LaksnF8tjoSPVcvS2xDw'
        let signature = ecc.sign('otcexchange',priv)
        console.log(signature)
    }
    f()

