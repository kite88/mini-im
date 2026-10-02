#!/usr/bin/env node
/**
 * Bootstrap CSS 精简脚本
 * ==========================================================================
 * 背景：项目没有前端构建流程，bootstrap.min.css 是 vendor 进来的完整包，
 *      但页面实际只用到网格 / 按钮 / 表单等一小部分。
 *
 * 做法：采集项目里所有可能出现在 DOM 上的 class 名，然后按「选择器分支」裁剪：
 *   1. 分支里不含 class（元素选择器、属性选择器、box-sizing 等）—— 保留
 *   2. 分支里的 class 全部在「已使用集合」—— 保留
 *   3. 分支引用了没人用的 class —— 删除该分支
 *   4. 一条规则的所有分支都被删掉时，整条规则删除；否则保留剩余分支
 *
 * 第 4 点是关键：`.checkbox,.radio{}` 这种共用声明必须拆开处理，
 * 否则 `.radio` 没人用会连 `.checkbox` 一起被误删。
 *
 * 安全边界：被删掉的分支一定含有当前 DOM 不存在的 class，结构上不可能命中，
 *          因此对现有页面而言裁剪前后渲染等价。@media / @supports 递归处理，
 *          @keyframes 等整块保留，授权注释保留。
 *
 * 用法：
 *   node tools/trim-bootstrap.cjs            # 空跑，只打印统计与校验
 *   node tools/trim-bootstrap.cjs --write    # 校验通过后写入
 *
 * 重新下载完整版 bootstrap 后需要再跑一次本脚本。
 */
const fs = require('fs');
const path = require('path');

const ROOT = path.resolve(__dirname, '..');
const CSS_FILE = path.join(ROOT, 'web/assets/lib/bootstrap-3.3.7/css/bootstrap.min.css');

// 参与 class 采集的文件
const SOURCE_FILES = [
    'web/index.html',
    'web/login.html',
    'web/conf.js',
    'web/assets/js/profile.js',
    'web/assets/js/cache.data.js',
    'web/assets/js/util/cookie.util.js',
    'web/assets/js/util/crypt.util.js',
    'web/assets/js/util/date.util.js',
    'web/assets/css/1.0.0/common.css',
    'web/assets/css/1.0.0/login.css',
    'web/assets/css/1.0.0/style.css',
];

const TOKEN = /^[A-Za-z][\w-]*$/;
const CLASS_IN_SELECTOR = /\.(-?[_A-Za-z][\w-]*)/g;
const NESTED_AT = ['@media', '@supports', '@document', '@layer', '@container'];

/** 采集「可能出现在 DOM 上的 class 名」（宁可多采，多采只会少删） */
function collectUsedClasses() {
    const used = new Set();
    const add = (value) => {
        for (const token of String(value).split(/[\s,;>+~()[\]{}:"'`=.|]+/)) {
            if (TOKEN.test(token)) used.add(token);
        }
    };

    for (const rel of SOURCE_FILES) {
        const file = path.join(ROOT, rel);
        if (!fs.existsSync(file)) continue;
        const text = fs.readFileSync(file, 'utf8');

        // a) 静态 class 属性
        for (const m of text.matchAll(/class\s*=\s*"([^"]*)"/g)) add(m[1]);
        // b) JS 里的字符串字面量，覆盖 v-bind:class、拼出来的 class 等
        if (rel.endsWith('.js')) {
            for (const m of text.matchAll(/'([^'\n]*)'|"([^"\n]*)"/g)) add(m[1] || m[2]);
        }
    }
    return used;
}

/** 抽出注释，保留 /*! 开头的授权注释 */
function extractComments(css) {
    const licenses = [];
    const stripped = css.replace(/\/\*[\s\S]*?\*\//g, (comment) => {
        if (comment.startsWith('/*!')) licenses.push(comment);
        return '';
    });
    return { stripped, licenses };
}

/** 按花括号切分出 { 选择器, 声明体 } 规则列表（嵌套块交给调用方递归） */
function parseRules(css) {
    const rules = [];
    let buf = '';
    let i = 0;

    while (i < css.length) {
        if (css[i] === '{') {
            let depth = 1;
            let j = i + 1;
            while (j < css.length && depth > 0) {
                if (css[j] === '{') depth++;
                else if (css[j] === '}') depth--;
                j++;
            }
            rules.push({ prelude: buf.trim(), body: css.slice(i + 1, j - 1) });
            buf = '';
            i = j;
            continue;
        }
        buf += css[i];
        i++;
    }
    return { rules, trailing: buf.trim() };
}

/** 按顶层逗号拆分选择器列表，忽略 :not(...) / [attr=","] / 引号里的逗号 */
function splitSelectorList(selector) {
    const parts = [];
    let buf = '';
    let depth = 0;
    let quote = null;

    for (let i = 0; i < selector.length; i++) {
        const ch = selector[i];
        if (quote) {
            buf += ch;
            if (ch === quote && selector[i - 1] !== '\\') quote = null;
            continue;
        }
        if (ch === '"' || ch === "'") {
            quote = ch;
            buf += ch;
            continue;
        }
        if (ch === '(' || ch === '[') depth++;
        else if (ch === ')' || ch === ']') depth--;

        if (ch === ',' && depth === 0) {
            parts.push(buf.trim());
            buf = '';
            continue;
        }
        buf += ch;
    }
    if (buf.trim()) parts.push(buf.trim());
    return parts;
}

const classesOf = (selector) =>
    [...selector.matchAll(CLASS_IN_SELECTOR)].map((m) => m[1]);

function filterRules(css, used, stats) {
    const { rules, trailing } = parseRules(css);
    let out = '';

    if (trailing) {
        // @charset / @import 之类块外文本：原样保留并计数
        out += trailing;
        stats.trailingText++;
    }

    for (const { prelude, body } of rules) {
        if (!prelude) continue;

        if (prelude.startsWith('@')) {
            if (NESTED_AT.some((prefix) => prelude.startsWith(prefix))) {
                const inner = filterRules(body, used, stats);
                if (inner) out += `${prelude}{${inner}}`;
                else stats.droppedBlocks++;
            } else {
                // @font-face / @keyframes / @page 等：整块保留
                out += `${prelude}{${body}}`;
                stats.keptBlocks++;
            }
            continue;
        }

        const parts = splitSelectorList(prelude);
        const kept = parts.filter((part) => {
            const classes = classesOf(part);
            return classes.length === 0 || classes.every((c) => used.has(c));
        });

        stats.totalParts += parts.length;
        stats.keptParts += kept.length;
        stats.droppedParts += parts.length - kept.length;

        if (kept.length === 0) {
            stats.droppedRules++;
            continue;
        }
        out += `${kept.join(',')}{${body}}`;
        stats.keptRules++;
        stats.emittedBodies.push(body);
    }

    return out;
}

function main() {
    const write = process.argv.includes('--write');
    const used = collectUsedClasses();
    const original = fs.readFileSync(CSS_FILE, 'utf8');
    const { stripped, licenses } = extractComments(original);

    const stats = {
        keptRules: 0, droppedRules: 0, keptBlocks: 0, droppedBlocks: 0,
        keptParts: 0, droppedParts: 0, totalParts: 0, trailingText: 0, emittedBodies: [],
    };
    const body = filterRules(stripped, used, stats);
    const output = `${licenses.join('\n')}\n${body}\n`;

    console.log(`已采集 class 名：${used.size} 个`);
    console.log(`规则：保留 ${stats.keptRules} 条 / 删除 ${stats.droppedRules} 条`);
    console.log(`选择器分支：保留 ${stats.keptParts} / 删除 ${stats.droppedParts} / 合计 ${stats.totalParts}`);
    console.log(`嵌套块：保留 ${stats.keptBlocks} / 空置删除 ${stats.droppedBlocks}；块外文本 ${stats.trailingText}`);
    console.log(`授权注释：保留 ${licenses.length} 条`);
    console.log(`体积：${original.length} -> ${output.length} bytes（省 ${(100 - (output.length / original.length) * 100).toFixed(1)}%）`);

    const checks = [];

    // 校验 1：拆分守恒
    checks.push([
        '分支总数守恒',
        stats.keptParts + stats.droppedParts === stats.totalParts,
        `${stats.keptParts}+${stats.droppedParts}=${stats.totalParts}`,
    ]);

    // 校验 2：输出里每条分支的 class 都在使用集合内
    const leftovers = new Set();
    const outputNoComment = output.replace(/\/\*[\s\S]*?\*\//g, '');
    for (const part of splitSelectorList(
        [...outputNoComment.matchAll(/([^{}]+)\{/g)].map((m) => m[1]).join(',')
    )) {
        for (const c of classesOf(part)) if (!used.has(c)) leftovers.add(c);
    }
    checks.push(['无残留未使用 class', leftovers.size === 0, [...leftovers].join(',') || '-']);

    // 校验 3：每个被删除的分支都必须含有至少一个未使用的 class
    const droppedBad = [];
    const droppedParts = [];
    (function walk(css) {
        const { rules } = parseRules(css);
        for (const { prelude, body: b } of rules) {
            if (!prelude) continue;
            if (prelude.startsWith('@')) {
                if (NESTED_AT.some((p) => prelude.startsWith(p))) walk(b);
                continue;
            }
            for (const part of splitSelectorList(prelude)) {
                const classes = classesOf(part);
                const keep = classes.length === 0 || classes.every((c) => used.has(c));
                if (!keep) droppedParts.push(part);
            }
        }
    })(stripped);
    for (const part of droppedParts) {
        if (!classesOf(part).some((c) => !used.has(c))) droppedBad.push(part);
    }
    checks.push(['被删分支均含未使用 class', droppedBad.length === 0, `${droppedBad.length} 个异常`]);

    // 校验 4：保留分支的声明体必须与原文件逐字节一致（防止声明被改动）
    const originalBodies = new Set([...stripped.matchAll(/\{([^{}]*)\}/g)].map((m) => m[1]));
    const mutated = stats.emittedBodies.filter((b) => !originalBodies.has(b));
    checks.push(['声明体未被改动', mutated.length === 0, `${mutated.length} 处不一致`]);

    // 校验 5：括号配平
    const balance = (output.match(/\{/g) || []).length - (output.match(/\}/g) || []).length;
    checks.push(['括号配平', balance === 0, String(balance)]);

    let ok = true;
    for (const [name, pass, detail] of checks) {
        if (!pass) ok = false;
        console.log(`${pass ? '[通过]' : '[失败]'} ${name}${detail && detail !== '-' ? `（${detail}）` : ''}`);
    }

    if (write && ok) {
        fs.writeFileSync(CSS_FILE, output);
        console.log('已写入。');
    } else if (write) {
        console.error('存在校验失败项，未写入。');
        process.exit(1);
    }
}

main();
