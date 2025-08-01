document.addEventListener('DOMContentLoaded', function () {
  // Check if geolocation has already been processed in this session
  if (!sessionStorage.getItem('geolocation_processed')) {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        function (position) {
          const lat = position.coords.latitude;
          const lon = position.coords.longitude;
          sessionStorage.setItem('geolocation_processed', 'true'); // Set flag
          window.location.href = `/weather/latlon?lat=${lat}&lon=${lon}`;
        },
        function (error) {
          if (error.code === error.PERMISSION_DENIED) {
            sessionStorage.setItem('geolocation_processed', 'true'); // Set flag even on permission denied
            window.location.href = '/ip';
          } else {
            // Handle other geolocation errors, e.g., position unavailable
            sessionStorage.setItem('geolocation_processed', 'true');
            window.location.href = '/ip'; // Fallback to IP-based lookup
          }
        }
      );
    } else {
      alert('Geolocation is not supported by this browser.');
      sessionStorage.setItem('geolocation_processed', 'true'); // Set flag
      window.location.href = '/ip';
    }
  }

  // Initialize GSAP animations
  gsap.registerPlugin(ScrollTrigger);

  // Animate all elements with .animate-in class
  gsap.utils.toArray('.animate-in').forEach((element, index) => {
    gsap.fromTo(element,
      { opacity: 0, y: 30 },
      {
        opacity: 1,
        y: 0,
        duration: 0.6,
        delay: index * 0.1,
        scrollTrigger: {
          trigger: element,
          start: "top 90%",
          toggleActions: "play none none none"
        }
      }
    );
  });

  // Add hover animations to weather details
  gsap.utils.toArray('.weather-detail').forEach(element => {
    element.addEventListener('mouseenter', () => {
      gsap.to(element, {
        y: -5,
        duration: 0.3,
        boxShadow: "0 10px 25px rgba(0,0,0,0.2)"
      });
    });

    element.addEventListener('mouseleave', () => {
      gsap.to(element, {
        y: 0,
        duration: 0.3,
        boxShadow: "0 4px 6px rgba(0,0,0,0.1)"
      });
    });
  });

  // Floating animation for weather icon
  gsap.to('.floating', {
    y: -10,
    duration: 2,
    repeat: -1,
    yoyo: true,
    ease: "sine.inOut"
  });

  // Sun position animation
  const sunPosition = document.querySelector('.sun-position');
  if (sunPosition) {
    gsap.fromTo(sunPosition,
      { left: "0%" },
      {
        left: sunPosition.style.left,
        duration: 2,
        ease: "power2.out"
      }
    );
  }

  // Background elements animation
  gsap.to('.absolute > div', {
    y: (i) => i % 2 === 0 ? -10 : 10,
    duration: 3,
    repeat: -1,
    yoyo: true,
    ease: "sine.inOut",
    stagger: 0.5
  });

  // Search functionality
  const searchInput = document.getElementById('search-input');
  const searchButton = document.getElementById('search-button');

  if (searchInput && searchButton) {
    const performSearch = () => {
      const query = searchInput.value.trim();
      if (query) {
        window.location.href = `/weather?query=${encodeURIComponent(query)}`;
      }
    };

    searchButton.addEventListener('click', performSearch);
    searchInput.addEventListener('keypress', (event) => {
      if (event.key === 'Enter') {
        performSearch();
      }
    });
  }
});