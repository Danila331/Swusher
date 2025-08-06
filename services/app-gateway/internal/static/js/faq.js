// FAQ Accordion
document.addEventListener('DOMContentLoaded', function() {
    const faqItems = document.querySelectorAll('.faq-item');
    
    faqItems.forEach(item => {
        const question = item.querySelector('.faq-question');
        const answer = item.querySelector('.faq-answer');
        const toggle = item.querySelector('.faq-toggle i');
        
        // Set initial state
        answer.style.maxHeight = '0';
        
        question.addEventListener('click', () => {
            // Close all other items
            faqItems.forEach(otherItem => {
                if (otherItem !== item && otherItem.classList.contains('active')) {
                    otherItem.classList.remove('active');
                    const otherAnswer = otherItem.querySelector('.faq-answer');
                    const otherToggle = otherItem.querySelector('.faq-toggle i');
                    otherAnswer.style.maxHeight = '0';
                    otherToggle.style.transform = 'rotate(0deg)';
                }
            });
            
            // Toggle current item
            const isActive = item.classList.contains('active');
            
            if (!isActive) {
                item.classList.add('active');
                answer.style.maxHeight = answer.scrollHeight + 'px';
                toggle.style.transform = 'rotate(180deg)';
            } else {
                item.classList.remove('active');
                answer.style.maxHeight = '0';
                toggle.style.transform = 'rotate(0deg)';
            }
        });
    });
}); 