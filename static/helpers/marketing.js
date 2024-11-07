const clickIdKeys = ["gclid", "gbraid", "wbraid", "msclkid"];

export class MarketingHelper {
    constructor() {
        this.user = JSON.parse(localStorage.getItem("user")) || {};
        this.landingPage = new URL(this.user.landingPage || window.location.href);
        this.language = navigator.language || navigator.userLanguage;
        this.marketingParams = Object.fromEntries(this.landingPage.searchParams);

        this.clickId = null;
        this.facebookClickId = null;
        this.longitude = null;
        this.latitude = null;

        // Get Click ID
        if (this.isPaid(this.landingPage.searchParams)) {
            this.clickId = this.getClickId(this.landingPage.searchParams);
        };

        // Get FB Click ID
        const fbClickId = this.landingPage.searchParams.get("fbclid");
        if (fbClickId) this.facebookClickId = fbClickId;

        this.data = new FormData();
    }

    getMarketingData() {
        this.data.append("click_id", this.clickId);
        if (this.facebookClickId) this.data.append("facebook_click_id", this.facebookClickId);
        this.data.append("landing_page", this.user.landingPage);
        this.data.append("referrer", this.user.referrer);
        this.data.append("language", this.language);


        // Append source, medium, and channel based on URL or referrer
        const source = this.landingPage.searchParams.get("source") || this.getSource(this.user.referrer);
        const medium = this.landingPage.searchParams.get("medium") || this.getMedium(this.user.referrer, this.landingPage.searchParams);
        const channel = this.landingPage.searchParams.get("channel") || this.getChannel(this.user.referrer);

        if (source) this.data.append("source", source);
        if (medium) this.data.append("medium", medium);
        if (channel) this.data.append("channel", channel);

        // Handle geolocation (conditionally append if available)
        if (this.longitude) this.data.append("longitude", this.longitude);
        if (this.latitude) this.data.append("latitude", this.latitude);

        // Append all marketing parameters and form values
        Object.entries(this.marketingParams).forEach(([key, value]) => value && this.data.append(key, value));

        return this.data;
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