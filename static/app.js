// Global variables
let startTime;

// Set example message in textarea
function setExample(message) {
    const messageInput = document.getElementById('customerMessage');
    if (messageInput) {
        messageInput.value = message;
    }
}

// Reset all flow steps to initial state
function resetFlow() {
    const steps = document.querySelectorAll('.flow-step');
    steps.forEach(function(step) {
        step.classList.remove('active', 'completed');
    });
    
    const resultPanel = document.getElementById('resultPanel');
    if (resultPanel) {
        resultPanel.classList.remove('show');
    }
}

// Activate a specific step in the flow
function activateStep(stepNumber) {
    const step = document.getElementById('step' + stepNumber);
    if (step) {
        step.classList.add('active');
        step.classList.remove('completed');
    }
}

// Mark a step as completed
function completeStep(stepNumber) {
    const step = document.getElementById('step' + stepNumber);
    if (step) {
        step.classList.remove('active');
        step.classList.add('completed');
    }
}

// Main function to classify customer message
async function classifyMessage() {
    const messageInput = document.getElementById('customerMessage');
    const classifyBtn = document.getElementById('classifyBtn');
    
    if (!messageInput || !classifyBtn) {
        console.error('Required elements not found');
        return;
    }
    
    const message = messageInput.value.trim();
    
    if (!message) {
        alert('Please enter a customer message');
        return;
    }
    
    // Reset and start the flow
    resetFlow();
    startTime = Date.now();
    
    // Disable button and show loading
    classifyBtn.disabled = true;
    classifyBtn.innerHTML = '<div class="loading"></div> Processing...';
    
    try {
        // Step 1: Message Received
        activateStep(1);
        await sleep(800);
        completeStep(1);
        
        // Step 2: AI Analysis
        activateStep(2);
        
        // Make API call
        const response = await fetch('/api/classify', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                customer_message: message
            })
        });
        
        if (!response.ok) {
            throw new Error('HTTP error! status: ' + response.status);
        }
        
        const data = await response.json();
        
        await sleep(1000);
        completeStep(2);
        
        // Step 3: Intent Classification
        activateStep(3);
        await sleep(800);
        completeStep(3);
        
        // Step 4: Agent Assignment
        activateStep(4);
        await sleep(800);
        completeStep(4);
        
        // Show results
        const endTime = Date.now();
        const processingTime = ((endTime - startTime) / 1000).toFixed(2);
        
        displayResults(data.intent, data.recommended_agent, processingTime);
        
    } catch (error) {
        console.error('Error:', error);
        alert('Classification failed: ' + error.message);
        resetFlow();
    } finally {
        // Re-enable button
        classifyBtn.disabled = false;
        classifyBtn.innerHTML = '🚀 Classify & Route Message';
    }
}

// Display the classification results
function displayResults(intent, agent, processingTime) {
    // Update result values
    const detectedIntentElement = document.getElementById('detectedIntent');
    const assignedAgentElement = document.getElementById('assignedAgent');
    const processingTimeElement = document.getElementById('processingTime');
    
    if (detectedIntentElement) {
        detectedIntentElement.textContent = intent.replace(/_/g, ' ').toUpperCase();
    }
    
    if (assignedAgentElement) {
        assignedAgentElement.textContent = agent.replace(/-/g, ' ').toUpperCase();
    }
    
    if (processingTimeElement) {
        processingTimeElement.textContent = processingTime + 's';
    }
    
    // Show result panel
    const resultPanel = document.getElementById('resultPanel');
    if (resultPanel) {
        resultPanel.classList.add('show');
    }
}

// Helper function for delays
function sleep(ms) {
    return new Promise(function(resolve) {
        setTimeout(resolve, ms);
    });
}

// Event listeners
document.addEventListener('DOMContentLoaded', function() {
    // Add enter key support for textarea
    const messageInput = document.getElementById('customerMessage');
    if (messageInput) {
        messageInput.addEventListener('keydown', function(e) {
            if (e.key === 'Enter' && e.ctrlKey) {
                classifyMessage();
            }
        });
    }
    
    // Add initial animation when page loads
    const container = document.querySelector('.container');
    if (container) {
        container.style.opacity = '0';
        container.style.transform = 'translateY(20px)';
        
        setTimeout(function() {
            container.style.transition = 'all 0.6s ease';
            container.style.opacity = '1';
            container.style.transform = 'translateY(0)';
        }, 100);
    }
});