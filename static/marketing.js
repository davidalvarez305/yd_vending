const qs = new URLSearchParams(window.location.search);
const clickIdKeys = ["gclid", "gbraid", "wbraid", "msclkid", "fbclid"];
let latitude = 0.0;
let longitude = 0.0;

document.addEventListener("DOMContentLoaded", getUserLocation());

function getHost(urlString) {
    let url;
    try {
        url = new URL(urlString);
    } catch (error) {
        return '';
    }

    let host = url.hostname.toLowerCase();

    if (host.startsWith('www.')) {
        host = host.slice(4);
    }

    const parts = host.split('.');

    // Handle cases like ftp.google.com or ads.google.com
    if (parts.length > 2 && !['com', 'net', 'org', 'edu', 'gov', 'mil', 'int'].includes(parts[parts.length - 2])) {
        host = parts.slice(-3).join('.');
    } else if (parts.length > 1) {
        // Check if the last part is a two-letter country code
        const lastPart = parts[parts.length - 1];
        if (lastPart.length === 2 && lastPart !== 'co') { // 'co' is a special case like .co.uk, .co.in, etc.
            host = parts.slice(-3).join('.');
        } else {
            host = parts.slice(-2).join('.');
        }
    }

    return host;
}

function getClickId() {
  for (const key of clickIdKeys) {
    if (qs.has(key)) return qs.get(key);
  }
}

function checkClickId() {
  for (const key of clickIdKeys) {
    if (qs.has(key)) {
      qs.delete(key);
      break;
    }
  }

  return qs.has("click_id");
}

function getMedium(referrer) {
  // No referrer means the user accessed the website directly
  if (referrer.length === 0) return "direct";

  // Non-empty referrer and no querystring === organic
  if (qs.size === 0) return "organic";

  // Paid ads
  if (checkClickId()) return "paid";

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
        { domain: "duckduckgo"},
        { domain: "yandex" },
        { domain: "baidu" },
        { domain: "naver" },
        { domain: "ask.com" },
        { domain: "adsensecustomsearchads" },
        { domain: "aol" },
        { domain: "brave" }
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
        "myspace"
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
        "vzaar"
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

function handleCTAClick(e) {
  const language = navigator.language || navigator.userLanguage;

  const buttonName = e.target.getAttribute("name");

  // Get user variables from browser
  var user = JSON.parse(localStorage.getItem("user")) || {};

  if (Object.keys(data).length === 0) {
    user.landingPage = window.location.href;
    user.referrer = document.referrer;
    localStorage.setItem("user", JSON.stringify(user));
  }

  qs.set("landing_page", user.landingPage);
  qs.set("referrer", user.referrer);
  qs.set("source", qs.get('source') ?? getHost(user.referrer)); // google.com || facebook.com || youtube.com
  qs.set("medium", qs.get('medium') ?? getMedium(user.referrer)); // organic || paid || direct
  qs.set("channel", qs.get('channel') ?? getChannel(user.referrer)); // search || social || video
  qs.set("button_clicked", buttonName);
  qs.set("longitude", longitude);
  qs.set("latitude", latitude);
  qs.set("language", language);
  if (checkClickId()) qs.set("click_id", getClickId());

  const currentDomain = new URL(window.location.origin + "/quote");

  currentDomain.search = qs.toString();

  window.location.replace(currentDomain.href);
}

function applyButtonlogic() {
  let quoteButtons = document.querySelectorAll(".quoteButton");

  quoteButtons.forEach((button) => {
    let children = button.children;
    Array.from(children).forEach((child) => {
      child.setAttribute("name", button.name);
    });

    button.addEventListener("click", handleCTAClick);
  });
};

applyButtonlogic();
