async function getTwin() {
    let token_id = document.getElementById('token_id').value
    const url = 'http://localhost:6660/twin?token_id='
    fetch(url+token_id)
    .then(response => response.json())
    .then(data => {
        document.getElementById('source_name').innerText = data.Source.name
        document.getElementById('source_img').src = data.Source.image
        document.getElementById('source_owner').innerText = data.Source.owner
        document.getElementById('source_hyperlink').innerText = 'Opensea.io/.../' + data.Source.id
        document.getElementById('source_hyperlink').href = data.Source.hyperlink

        document.getElementById('twin_name').innerText = data.Twin.name
        document.getElementById('twin_img').src = data.Twin.image
        document.getElementById('twin_owner').innerText = data.Twin.owner
        document.getElementById('twin_hyperlink').innerText = 'Opensea.io/.../' + data.Twin.id
        document.getElementById('twin_hyperlink').href = data.Twin.hyperlink
    })
    .catch(err => console.error(err))
}