const form = document.getElementById('feedbackForm');
const responseMsg = document.getElementById('responseMessage');

// Handle Form Submission
form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value,
        subject: document.getElementById('subject').value,
        message: document.getElementById('message').value
    };

    try {
        const response = await fetch('/api/feedback', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            const data = await response.json();
            responseMsg.style.color = 'green';
            responseMsg.textContent = `Success! Feedback ID ${data.id} created.`;
            form.reset();
        } else {
            throw new Error('Failed to submit');
        }
    } catch (error) {
        responseMsg.style.color = 'red';
        responseMsg.textContent = 'Error submitting form.';
        console.error(error);
    }
});