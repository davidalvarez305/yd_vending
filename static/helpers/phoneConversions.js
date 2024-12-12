const phoneNumber = document.getElementById("phoneNumberCTA");

phoneNumber.addEventListener("click", () => handlePhoneNumberClick())

function handlePhoneNumberClick() {
    // Report FB Conversion
    if (fbq) fbq("track", "{{ .LeadEventName }}");

    // Report LinkedIn
    if (window.lintrk) window.lintrk('track', { conversion_id: "{{ .LinkedInEventID }}" });

    // Report Microsoft Ads Conversion
    window.uetq = window.uetq || [];

    window.uetq.push('event', '{{ .QuoteEventName }}', {});
}