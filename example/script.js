// ES5.1-compatible configuration script for clash-converter.

var RULE_BASE = 'https://cdn.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@meta/geo/';

// [policy group, source type, source name]
// Order matters: specific services are matched before CN and global fallbacks.
var RULESET_DEFINITIONS = [
    ['Ads', 'geosite', 'category-ads-all'],
    ['Private', 'geosite', 'private'],
    ['Private', 'geoip', 'private'],

    ['OpenAI', 'geosite', 'openai'],
    ['Anthropic', 'geosite', 'anthropic'],
    ['AI Chat', 'geosite', 'category-ai-chat-!cn'],
    ['YouTube', 'geosite', 'youtube'],
    ['Google', 'geosite', 'google'],
    ['Google', 'geoip', 'google'],
    ['OneDrive', 'geosite', 'onedrive'],
    ['Microsoft', 'geosite', 'microsoft'],
    ['Apple TV+', 'geosite', 'apple-tvplus'],
    ['iCloud', 'geosite', 'icloud'],
    ['Apple', 'geosite', 'apple'],

    ['Telegram', 'geosite', 'telegram'],
    ['Telegram', 'geoip', 'telegram'],
    ['Twitter/X', 'geosite', 'twitter'],
    ['Twitter/X', 'geoip', 'twitter'],
    ['Meta', 'geosite', 'facebook'],
    ['Meta', 'geosite', 'instagram'],
    ['Meta', 'geosite', 'whatsapp'],
    ['Meta', 'geoip', 'facebook'],
    ['Discord', 'geosite', 'discord'],
    ['TikTok', 'geosite', 'tiktok'],
    ['LINE', 'geosite', 'line'],
    ['Reddit', 'geosite', 'reddit'],
    ['LinkedIn', 'geosite', 'linkedin'],

    ['Netflix', 'geosite', 'netflix'],
    ['Netflix', 'geoip', 'netflix'],
    ['Disney+', 'geosite', 'disney'],
    ['HBO', 'geosite', 'hbo'],
    ['Hulu', 'geosite', 'hulu'],
    ['Prime Video', 'geosite', 'primevideo'],
    ['Spotify', 'geosite', 'spotify'],
    ['Twitch', 'geosite', 'twitch'],
    ['Bahamut', 'geosite', 'bahamut'],
    ['BiliIntl', 'geosite', 'biliintl'],

    ['Steam', 'geosite', 'steam'],
    ['Epic', 'geosite', 'epicgames'],
    ['EA', 'geosite', 'ea'],
    ['Ubisoft', 'geosite', 'ubisoft'],
    ['Blizzard', 'geosite', 'blizzard'],
    ['GOG', 'geosite', 'gog'],
    ['Riot', 'geosite', 'riot'],
    ['PlayStation', 'geosite', 'playstation'],
    ['Xbox', 'geosite', 'xbox'],
    ['Nintendo', 'geosite', 'nintendo'],

    ['GitHub', 'geosite', 'github'],
    ['GitLab', 'geosite', 'gitlab'],
    ['Atlassian', 'geosite', 'atlassian'],
    ['Docker', 'geosite', 'docker'],
    ['npm', 'geosite', 'npmjs'],
    ['JetBrains', 'geosite', 'jetbrains'],
    ['StackExchange', 'geosite', 'stackexchange'],

    ['AWS', 'geosite', 'aws'],
    ['Azure', 'geosite', 'azure'],
    ['Cloudflare', 'geosite', 'cloudflare'],
    ['Cloudflare', 'geoip', 'cloudflare'],
    ['DigitalOcean', 'geosite', 'digitalocean'],
    ['Vercel', 'geosite', 'vercel'],
    ['Netlify', 'geosite', 'netlify'],

    ['PayPal', 'geosite', 'paypal'],
    ['Stripe', 'geosite', 'stripe'],
    ['Wise', 'geosite', 'wise'],
    ['Binance', 'geosite', 'binance'],

    ['Academic', 'geosite', 'coursera'],
    ['Academic', 'geosite', 'udemy'],
    ['Academic', 'geosite', 'edx'],
    ['Academic', 'geosite', 'khanacademy'],
    ['Academic', 'geosite', 'google-scholar'],
    ['Academic', 'geosite', 'category-scholar-!cn'],
    ['News', 'geosite', 'bbc'],
    ['News', 'geosite', 'cnn'],
    ['News', 'geosite', 'nytimes'],
    ['News', 'geosite', 'wsj'],
    ['News', 'geosite', 'bloomberg'],
    ['Shopping', 'geosite', 'amazon'],
    ['Shopping', 'geosite', 'ebay'],

    ['CN', 'geosite', 'geolocation-cn'],
    ['CN', 'geoip', 'cn'],
    ['PROXY', 'geosite', 'geolocation-!cn']
];

var POLICY_GROUPS = [
    'OpenAI', 'Anthropic', 'AI Chat',
    'YouTube', 'Google', 'OneDrive', 'Microsoft', 'Apple TV+', 'iCloud', 'Apple',
    'Telegram', 'Twitter/X', 'Meta', 'Discord', 'TikTok', 'LINE', 'Reddit', 'LinkedIn',
    'Netflix', 'Disney+', 'HBO', 'Hulu', 'Prime Video', 'Spotify', 'Twitch', 'Bahamut', 'BiliIntl',
    'Steam', 'Epic', 'EA', 'Ubisoft', 'Blizzard', 'GOG', 'Riot', 'PlayStation', 'Xbox', 'Nintendo',
    'GitHub', 'GitLab', 'Atlassian', 'Docker', 'npm', 'JetBrains', 'StackExchange',
    'AWS', 'Azure', 'Cloudflare', 'DigitalOcean', 'Vercel', 'Netlify',
    'PayPal', 'Stripe', 'Wise', 'Binance',
    'Academic', 'News', 'Shopping'
];

var INFO_NODE_PATTERN = /流量|到期|剩余|套餐|expire|traffic|quota|bandwidth/i;
var TEST_URL = 'https://www.gstatic.com/generate_204';

function rulesets(register) {
    var i;
    var definition;
    var behavior;
    var url;

    for (i = 0; i < RULESET_DEFINITIONS.length; i += 1) {
        definition = RULESET_DEFINITIONS[i];
        behavior = definition[1] === 'geoip' ? 'ipcidr' : 'domain';
        url = RULE_BASE + definition[1] + '/' + definition[2] + '.list';
        register(definition[0], url, behavior);
    }
}

function unique(values) {
    var result = [];
    var seen = {};
    var i;
    var value;

    for (i = 0; i < values.length; i += 1) {
        value = values[i];
        if (!seen[value]) {
            seen[value] = true;
            result.push(value);
        }
    }
    return result;
}

function selectGroup(name, proxies) {
    return {
        name: name,
        type: 'select',
        proxies: unique(proxies)
    };
}

function serviceCandidates(proxyNames) {
    return ['PROXY', 'DIRECT', 'AUTO', 'FALLBACK'].concat(proxyNames);
}

function buildGroups(proxyNames) {
    var groups = [
        selectGroup('PROXY', ['AUTO', 'FALLBACK', 'DIRECT'].concat(proxyNames)),
        {
            name: 'AUTO',
            type: 'url-test',
            url: TEST_URL,
            interval: 300,
            tolerance: 50,
            lazy: true,
            proxies: proxyNames
        },
        {
            name: 'FALLBACK',
            type: 'fallback',
            url: TEST_URL,
            interval: 300,
            lazy: true,
            proxies: proxyNames
        },
        selectGroup('Ads', ['REJECT', 'DIRECT', 'PROXY']),
        selectGroup('Private', ['DIRECT', 'PROXY'].concat(proxyNames)),
        selectGroup('CN', ['DIRECT', 'PROXY', 'AUTO', 'FALLBACK'].concat(proxyNames))
    ];
    var candidates = serviceCandidates(proxyNames);
    var i;

    for (i = 0; i < POLICY_GROUPS.length; i += 1) {
        groups.push(selectGroup(POLICY_GROUPS[i], candidates));
    }
    groups.push(selectGroup('Final', candidates));
    return groups;
}

function buildConfig(config, legacyRelay) {
    var proxies = Array.isArray(config['proxies']) ? config['proxies'] : [];
    var filteredProxies = [];
    var proxyNames = [];
    var i;
    var proxy;
    var name;
    var originalName;

    for (i = 0; i < proxies.length; i += 1) {
        proxy = proxies[i];
        name = proxy && typeof proxy.name === 'string' ? proxy.name.replace(/^\s+|\s+$/g, '') : '';
        originalName = name.replace(/^\[[^\]]+\]\s*/, '');
        if (!name || INFO_NODE_PATTERN.test(originalName)) {
            continue;
        }
        proxy.name = name;
        filteredProxies.push(proxy);
        proxyNames.push(name);
    }

    config['proxies'] = filteredProxies;
    config['proxy-groups'] = buildGroups(unique(proxyNames));
    config['rules'] = (Array.isArray(config['rules']) ? config['rules'] : []).concat(['MATCH,Final']);
}
