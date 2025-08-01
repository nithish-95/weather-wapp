
document.addEventListener('DOMContentLoaded', function() {
  // Check if geolocation has already been processed in this session
  if (sessionStorage.getItem('geolocation_processed')) {
    // If already processed, do nothing or redirect to a default page if needed
    // For now, we'll just prevent re-execution of geolocation logic
    return;
  }

  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      function(position) {
        const lat = position.coords.latitude;
        const lon = position.coords.longitude;
        sessionStorage.setItem('geolocation_processed', 'true'); // Set flag
        window.location.href = `/weather/latlon?lat=${lat}&lon=${lon}`;
      },
      function(error) {
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
});