document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('contractForm');

    form.addEventListener('submit', function(event) {
        event.preventDefault();

        const clientName = document.getElementById('clientName').value;
        const clientEmail = document.getElementById('clientEmail').value;
        const paymentAmount = document.getElementById('paymentAmount').value;
        const requirements = document.getElementById('requirements').value;
        const description = document.getElementById('description').value;

        const formData = {
            clientName: clientName,
            clientEmail: clientEmail,
            paymentAmount: paymentAmount,
            requirements: requirements,
            description: description
        };

        fetch('/generate_contract', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        })
        .then(response => response.text())
        .then(data => {
            alert(data); // Show the response from the server
        })
        .catch((error) => {
            console.error('Error:', error);
        });
    });
});
