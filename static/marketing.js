const clickIdKeys = ["gclid", "gbraid", "wbraid", "msclkid", "fbclid"];
const form = document.getElementById("get-a-quote-form");
let latitude = 0.0;
let longitude = 0.0;

document.addEventListener("DOMContentLoaded", getUserLocation());

function getHost(urlString) {
  let url;
  try {
    url = new URL(urlString);
  } catch (error) {
    return "";
  }

  let host = url.hostname.toLowerCase();

  if (host.startsWith("www.")) {
    host = host.slice(4);
  }

  const parts = host.split(".");

  // Handle cases like ftp.google.com or ads.google.com
  if (
    parts.length > 2 &&
    !["com", "net", "org", "edu", "gov", "mil", "int"].includes(
      parts[parts.length - 2]
    )
  ) {
    host = parts.slice(-3).join(".");
  } else if (parts.length > 1) {
    // Check if the last part is a two-letter country code
    const lastPart = parts[parts.length - 1];
    if (lastPart.length === 2 && lastPart !== "co") {
      // 'co' is a special case like .co.uk, .co.in, etc.
      host = parts.slice(-3).join(".");
    } else {
      host = parts.slice(-2).join(".");
    }
  }

  return host;
}

function getClickId(qs) {
  for (const key of clickIdKeys) {
    if (qs.has(key)) return qs.get(key);
  }
}

function isPaid(qs) {
  for (const key of clickIdKeys) {
    if (qs.has(key)) {
      return true;
    }
  }

  return false;
}

function getMedium(referrer, qs) {
  // No referrer means the user accessed the website directly
  if (referrer.length === 0) return "direct";

  // Non-empty referrer and no querystring === organic
  if (qs.size === 0) return "organic";

  // Paid ads
  if (isPaid(qs)) return "paid";

  // Querystring + non-empty referrer and no click id === referral
  return "referral";
}

function getUserLocation() {
  const options = {
    enableHighAccuracy: true,
    timeout: 3000,
    maximumAge: 0,
  };

  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      function (position) {
        latitude = position.coords.latitude;
        longitude = position.coords.longitude;
      },
      function (error) {
        console.error("Error getting user location:", error.message);
      },
      options
    );
  } else {
    console.error("Geolocation is not supported by your browser.");
  }
}

function getChannel(referrerUrl) {
  const searchEngines = [
    { domain: "google" },
    { domain: "bing" },
    { domain: "yahoo" },
    { domain: "ecosia" },
    { domain: "duckduckgo" },
    { domain: "yandex" },
    { domain: "baidu" },
    { domain: "naver" },
    { domain: "ask.com" },
    { domain: "adsensecustomsearchads" },
    { domain: "aol" },
    { domain: "brave" },
  ];

  const majorSocialNetworks = [
    "facebook",
    "instagram",
    "twitter",
    "linkedin",
    "pinterest",
    "snapchat",
    "reddit",
    "whatsapp",
    "wechat",
    "telegram",
    "discord",
    "vkontakte",
    "weibo",
    "line",
    "kakaotalk",
    "qq",
    "viber",
    "telegram",
    "tumblr",
    "flickr",
    "meetup",
    "tagged",
    "badoo",
    "myspace",
  ];

  const majorVideoPlatforms = [
    "youtube",
    "tiktok",
    "vimeo",
    "dailymotion",
    "twitch",
    "bilibili",
    "youku",
    "rutube",
    "vine",
    "peertube",
    "ig tv",
    "veoh",
    "metacafe",
    "vudu",
    "vidyard",
    "rumble",
    "bit chute",
    "brightcove",
    "viddler",
    "vzaar",
  ];

  // Check search engines
  for (let engine of searchEngines) {
    if (referrerUrl.includes(engine.domain)) {
      return "search";
    }
  }

  // Check social networks
  for (let network of majorSocialNetworks) {
    if (referrerUrl.includes(network)) {
      return "social";
    }
  }

  // Check video platforms
  for (let platform of majorVideoPlatforms) {
    if (referrerUrl.includes(platform)) {
      return "video";
    }
  }

  return "other";
}

function applyButtonlogic() {
  let quoteButtons = document.querySelectorAll(".quoteButton");

  quoteButtons.forEach((button) => {
    let children = button.children;
    Array.from(children).forEach((child) => {
      child.setAttribute("name", button.name);
    });

    button.addEventListener("click", function () {
      // Bring form into focus
      form.scrollIntoView({ behavior: "smooth" });
      form.querySelector("input, textarea, select").focus();

      // Hide modal
      const modal = document.getElementById("modalOverlay");
      modal.style.display = "none";
    });
  });
}

applyButtonlogic();

form.addEventListener("submit", (e) => {
  e.preventDefault();

  const qs = new URLSearchParams(window.location.search);

  const language = navigator.language || navigator.userLanguage;

  const buttonName = e.target.getAttribute("name");

  // Get user variables from browser
  var user = JSON.parse(localStorage.getItem("user")) || {};

  const marketing = Object.fromEntries(qs);
  const data = new FormData(e.target);

  if (isPaid(qs)) {
    data.append("click_id", getClickId(qs));
  }

  data.append("landing_page", user.landingPage);
  data.append("referrer", user.referrer);
  data.append("source", qs.get("source") ?? getHost(user.referrer)); // google.com || facebook.com || youtube.com
  data.append("medium", qs.get("medium") ?? getMedium(user.referrer, qs)); // organic || paid || direct
  data.append("channel", qs.get("channel") ?? getChannel(user.referrer)); // search || social || video
  data.append("button_clicked", buttonName);
  data.append("longitude", longitude);
  data.append("latitude", latitude);
  data.append("language", language);

  for (const [key, value] of Object.entries(marketing)) {
    data.append(key, value);
  }

  const alertModal = document.getElementById("alertModal");
  fetch("/quote", {
    method: "POST",
    credentials: "include",
    body: data,
  })
    .then((response) => {
      if (response.ok) {
        const token = response.headers.get("X-Csrf-Token");

        if (token) {
          const csrf_token = document.getElementById("csrf_token");

          if (!csrf_token) return;

          csrf_token.value = token;
        }

        return response.text();
      } else {
        return response.text().then((err) => {
          throw new Error(err);
        });
      }
    })
    .then((html) => {
      alertModal.outerHTML = html;
      handleCloseAlertModal();

      form.reset();
    })
    .catch(console.error);
});
