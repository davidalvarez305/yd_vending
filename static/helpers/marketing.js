const clickIdKeys = ["gclid", "gbraid", "wbraid", "msclkid"];

export class MarketingHelper {
    constructor() {
        this.user = JSON.parse(localStorage.getItem("user")) || {};
        this.language = navigator.language || navigator.userLanguage;
        
        if (this.user.landingPage) {
            this.landingPage = new URL(this.user.landingPage);
        } else {
            this.landingPage = new URL(this.user.landingPage);
        }

        this.longitude = null;
        this.latitude = null;

        this.clickId = null;
        this.facebookClickId = null;
        this.facebookClientId = null;

        this.userAgent = navigator.userAgent;

        if (this.landingPage) {
            this.marketingParams = Object.fromEntries(this.landingPage.searchParams);

            // Get Click ID
            if (this.isPaid(this.landingPage.searchParams)) {
                this.clickId = this.getClickId(this.landingPage.searchParams);
            };

            // Get FB Click ID
            const fbClickId = this.landingPage.searchParams.get("fbclid");
            if (fbClickId) this.facebookClickId = fbClickId;
        }

        // Get FB Client ID
        const fbp = this.getCookie("_fbp");
        if (fbp) this.facebookClientId = fbp;

        this.data = new FormData();
    }

    populate() {
        if (this.clickId) this.data.set("click_id", this.clickId);
        if (this.facebookClickId) this.data.set("facebook_click_id", this.facebookClickId);
        if (this.facebookClientId) this.data.set("facebook_client_id", this.facebookClientId);

        if (this.landingPage) this.data.set("landing_page", this.landingPage);
        if (this.user.referrer) this.data.set("referrer", this.user.referrer);
        if (this.language) this.data.set("language", this.language);

        // Append source, medium, and channel based on URL or referrer
        const source = this.landingPage.searchParams.get("source") || this.getSource(this.user.referrer);
        const medium = this.landingPage.searchParams.get("medium") || this.getMedium(this.user.referrer, this.landingPage.searchParams);
        const channel = this.landingPage.searchParams.get("channel") || this.getChannel(this.user.referrer);

        if (source) this.data.set("source", source);
        if (medium) this.data.set("medium", medium);
        if (channel) this.data.set("channel", channel);

        // Handle geolocation (conditionally append if available)
        if (this.longitude) this.data.set("longitude", this.longitude);
        if (this.latitude) this.data.set("latitude", this.latitude);

        // Append all marketing parameters
        Object.entries(this.marketingParams).forEach(([key, value]) => {
            if (value) this.data.set(key, value);
        });
    }

    getCookie(name) {
        const cookies = document.cookie.split(';');

        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith(name + '=')) {
                return cookie.substring(name.length + 1);
            }
        }

        return null;
    }

    getSource(urlString) {
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

    getClickId(qs) {
        for (const key of clickIdKeys) {
            if (qs.has(key)) return qs.get(key);
        }
    }

    isPaid(qs) {
        for (const key of clickIdKeys) {
            if (qs.has(key)) {
                return true;
            }
        }

        return false;
    }

    getMedium(referrer, qs) {
        if (!referrer) return "direct";

        // No referrer means the user accessed the website directly
        if (referrer.length === 0) return "direct";

        // Non-empty referrer and no querystring === organic
        if (qs.size === 0) return "organic";

        // Paid ads
        if (this.isPaid(qs)) return "paid";

        // Querystring + non-empty referrer and no click id === referral
        return "referral";
    }

    getUserLocation() {
        const options = {
            enableHighAccuracy: true,
            timeout: 3000,
            maximumAge: 0,
        };

        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                function (position) {
                    this.latitude = position.coords.latitude;
                    this.longitude = position.coords.longitude;
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

    getChannel(referrerUrl) {
        if (!referrerUrl) return "other";

        const displayNetworks = ["googleads.g.doubleclick.net"];

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

        // Check display platforms
        for (let platform of displayNetworks) {
            if (referrerUrl.includes(platform)) {
                return "display";
            }
        }

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
}