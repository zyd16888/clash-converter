// ES5.1-compatible configuration script for clash-converter.

var RULE_BASE = 'https://cdn.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@meta/geo/';

var GROUP = {
    PROXY: '🚀 节点选择',
    AUTO: '⚡ 自动选择',
    FALLBACK: '🛟 故障转移',
    ADS: '🛑 广告拦截',
    AI: '🤖 AI 服务',
    YOUTUBE: '📹 YouTube',
    GOOGLE: '🔍 Google 服务',
    MICROSOFT: 'Ⓜ️ Microsoft 服务',
    APPLE: '🍏 Apple 服务',
    TELEGRAM: '📲 Telegram',
    SOCIAL: '💬 社交媒体',
    STREAMING: '🎬 国际流媒体',
    GAMING: '🎮 游戏平台',
    DEVELOPMENT: '🛠️ 开发工具',
    CLOUD_CN: '☁️ 国内云服务',
    CLOUD_GLOBAL: '🌐 国外云服务',
    FINANCE: '💳 金融支付',
    PRIVATE: '🏠 私有网络',
    CN: '🇨🇳 国内直连',
    FINAL: '🐟 漏网之鱼'
};

// [policy group, source type, source name]
// Order matters: specific services are matched before CN and global fallbacks.
var RULESET_DEFINITIONS = [
    [GROUP.ADS, 'geosite', 'category-ads-all'],
    [GROUP.PRIVATE, 'geosite', 'private'],
    [GROUP.PRIVATE, 'geoip', 'private'],

    [GROUP.CLOUD_CN, 'geosite', 'aliyun'],
    [GROUP.CLOUD_CN, 'geosite', 'huaweicloud'],
    [GROUP.CLOUD_CN, 'geosite', 'volcengine'],
    [GROUP.CLOUD_CN, 'geosite', 'ucloud'],
    [GROUP.CLOUD_CN, 'geosite', 'qiniu'],
    [GROUP.CLOUD_CN, 'geosite', 'aws-cn'],
    [GROUP.CLOUD_CN, 'geosite', 'azure@cn'],

    [GROUP.AI, 'geosite', 'google-gemini'],
    [GROUP.AI, 'geosite', 'google-deepmind'],
    [GROUP.AI, 'geosite', 'github-copilot'],
    [GROUP.AI, 'geosite', 'openai'],
    [GROUP.AI, 'geosite', 'anthropic'],
    [GROUP.AI, 'geosite', 'category-ai-chat-!cn'],
    [GROUP.YOUTUBE, 'geosite', 'youtube'],
    [GROUP.GOOGLE, 'geosite', 'google'],
    [GROUP.GOOGLE, 'geoip', 'google'],
    [GROUP.MICROSOFT, 'geosite', 'onedrive'],
    [GROUP.MICROSOFT, 'geosite', 'microsoft'],
    [GROUP.APPLE, 'geosite', 'apple-tvplus'],
    [GROUP.APPLE, 'geosite', 'icloud'],
    [GROUP.APPLE, 'geosite', 'apple'],

    [GROUP.TELEGRAM, 'geosite', 'telegram'],
    [GROUP.TELEGRAM, 'geoip', 'telegram'],
    [GROUP.SOCIAL, 'geosite', 'twitter'],
    [GROUP.SOCIAL, 'geoip', 'twitter'],
    [GROUP.SOCIAL, 'geosite', 'facebook'],
    [GROUP.SOCIAL, 'geosite', 'instagram'],
    [GROUP.SOCIAL, 'geosite', 'whatsapp'],
    [GROUP.SOCIAL, 'geoip', 'facebook'],
    [GROUP.SOCIAL, 'geosite', 'discord'],
    [GROUP.SOCIAL, 'geosite', 'tiktok'],
    [GROUP.SOCIAL, 'geosite', 'line'],
    [GROUP.SOCIAL, 'geosite', 'reddit'],
    [GROUP.SOCIAL, 'geosite', 'linkedin'],

    [GROUP.STREAMING, 'geosite', 'netflix'],
    [GROUP.STREAMING, 'geoip', 'netflix'],
    [GROUP.STREAMING, 'geosite', 'disney'],
    [GROUP.STREAMING, 'geosite', 'hbo'],
    [GROUP.STREAMING, 'geosite', 'hulu'],
    [GROUP.STREAMING, 'geosite', 'primevideo'],
    [GROUP.STREAMING, 'geosite', 'spotify'],
    [GROUP.STREAMING, 'geosite', 'twitch'],
    [GROUP.STREAMING, 'geosite', 'bahamut'],
    [GROUP.STREAMING, 'geosite', 'biliintl'],

    [GROUP.GAMING, 'geosite', 'steam'],
    [GROUP.GAMING, 'geosite', 'epicgames'],
    [GROUP.GAMING, 'geosite', 'ea'],
    [GROUP.GAMING, 'geosite', 'ubisoft'],
    [GROUP.GAMING, 'geosite', 'blizzard'],
    [GROUP.GAMING, 'geosite', 'gog'],
    [GROUP.GAMING, 'geosite', 'riot'],
    [GROUP.GAMING, 'geosite', 'playstation'],
    [GROUP.GAMING, 'geosite', 'xbox'],
    [GROUP.GAMING, 'geosite', 'nintendo'],

    [GROUP.DEVELOPMENT, 'geosite', 'github'],
    [GROUP.DEVELOPMENT, 'geosite', 'gitlab'],
    [GROUP.DEVELOPMENT, 'geosite', 'atlassian'],
    [GROUP.DEVELOPMENT, 'geosite', 'docker'],
    [GROUP.DEVELOPMENT, 'geosite', 'npmjs'],
    [GROUP.DEVELOPMENT, 'geosite', 'jetbrains'],
    [GROUP.DEVELOPMENT, 'geosite', 'stackexchange'],

    [GROUP.CLOUD_GLOBAL, 'geosite', 'aws'],
    [GROUP.CLOUD_GLOBAL, 'geosite', 'azure'],
    [GROUP.CLOUD_GLOBAL, 'geosite', 'cloudflare'],
    [GROUP.CLOUD_GLOBAL, 'geosite', 'digitalocean'],
    [GROUP.CLOUD_GLOBAL, 'geosite', 'vercel'],
    [GROUP.CLOUD_GLOBAL, 'geosite', 'netlify'],

    [GROUP.FINANCE, 'geosite', 'paypal'],
    [GROUP.FINANCE, 'geosite', 'stripe'],
    [GROUP.FINANCE, 'geosite', 'wise'],
    [GROUP.FINANCE, 'geosite', 'binance'],

    [GROUP.PROXY, 'geosite', 'coursera'],
    [GROUP.PROXY, 'geosite', 'udemy'],
    [GROUP.PROXY, 'geosite', 'edx'],
    [GROUP.PROXY, 'geosite', 'khanacademy'],
    [GROUP.PROXY, 'geosite', 'google-scholar'],
    [GROUP.PROXY, 'geosite', 'category-scholar-!cn'],
    [GROUP.PROXY, 'geosite', 'bbc'],
    [GROUP.PROXY, 'geosite', 'cnn'],
    [GROUP.PROXY, 'geosite', 'nytimes'],
    [GROUP.PROXY, 'geosite', 'wsj'],
    [GROUP.PROXY, 'geosite', 'bloomberg'],
    [GROUP.PROXY, 'geosite', 'amazon'],
    [GROUP.PROXY, 'geosite', 'ebay'],

    [GROUP.CN, 'geosite', 'geolocation-cn'],
    [GROUP.CN, 'geoip', 'cn'],
    [GROUP.PROXY, 'geosite', 'geolocation-!cn']
];

var POLICY_GROUPS = [
    GROUP.AI, GROUP.YOUTUBE, GROUP.GOOGLE, GROUP.MICROSOFT, GROUP.APPLE,
    GROUP.TELEGRAM, GROUP.SOCIAL, GROUP.STREAMING, GROUP.GAMING,
    GROUP.DEVELOPMENT, GROUP.CLOUD_CN, GROUP.CLOUD_GLOBAL, GROUP.FINANCE
];

// These providers have no dedicated cloud-only MetaCubeX source. Keep this
// list limited to cloud product domains so their consumer services stay out.
var DOMESTIC_CLOUD_DOMAIN_SUFFIXES = [
    'cloud.tencent.com',
    'tencentcloud.com',
    'tencentcloudapi.com',
    'myqcloud.com',
    'qcloud.com',
    'qcloudcdn.com',
    'qcloudcos.com',
    'qcloudimg.com',
    'cloudbase.net',
    'cloud.baidu.com',
    'baidubce.com',
    'bcebos.com',
    'bdcloudapi.com',
    'jdcloud.com',
    'jdcloud-api.com',
    'jdcloud-openapi.com',
    'jcloud.com',
    'jcloudcs.com',
    'ksyun.com',
    'ksyuncs.com'
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
    return [GROUP.PROXY, 'DIRECT', GROUP.AUTO, GROUP.FALLBACK].concat(proxyNames);
}

function directCandidates(proxyNames) {
    return ['DIRECT', GROUP.PROXY, GROUP.AUTO, GROUP.FALLBACK].concat(proxyNames);
}

function domesticCloudRules() {
    var rules = [];
    var i;

    for (i = 0; i < DOMESTIC_CLOUD_DOMAIN_SUFFIXES.length; i += 1) {
        rules.push('DOMAIN-SUFFIX,' + DOMESTIC_CLOUD_DOMAIN_SUFFIXES[i] + ',' + GROUP.CLOUD_CN);
    }
    return rules;
}

function buildGroups(proxyNames) {
    var groups = [
        selectGroup(GROUP.PROXY, [GROUP.AUTO, GROUP.FALLBACK, 'DIRECT'].concat(proxyNames)),
        {
            name: GROUP.AUTO,
            type: 'url-test',
            url: TEST_URL,
            interval: 300,
            tolerance: 50,
            lazy: true,
            proxies: proxyNames
        },
        {
            name: GROUP.FALLBACK,
            type: 'fallback',
            url: TEST_URL,
            interval: 300,
            lazy: true,
            proxies: proxyNames
        },
        selectGroup(GROUP.ADS, ['REJECT', 'DIRECT', GROUP.PROXY]),
        selectGroup(GROUP.PRIVATE, ['DIRECT', GROUP.PROXY].concat(proxyNames)),
        selectGroup(GROUP.CN, ['DIRECT', GROUP.PROXY, GROUP.AUTO, GROUP.FALLBACK].concat(proxyNames))
    ];
    var candidates = serviceCandidates(proxyNames);
    var domesticCandidates = directCandidates(proxyNames);
    var i;

    for (i = 0; i < POLICY_GROUPS.length; i += 1) {
        groups.push(selectGroup(
            POLICY_GROUPS[i],
            POLICY_GROUPS[i] === GROUP.CLOUD_CN ? domesticCandidates : candidates
        ));
    }
    groups.push(selectGroup(GROUP.FINAL, candidates));
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
    config['rules'] = domesticCloudRules()
        .concat(Array.isArray(config['rules']) ? config['rules'] : [])
        .concat(['MATCH,' + GROUP.FINAL]);
}
