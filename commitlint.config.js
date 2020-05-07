module.exports = {
    parserPreset: '@commitlint/config-conventional',
    rules: {
        'body-leading-blank': [1, 'always'],
        'header-max-length': [2, 'always', 50],
        'body-max-line-length': [2, 'always', 72],
        'footer-max-line-length': [2, 'always', 72],
        'footer-leading-blank': [1, 'always'],
        'scope-max-length': [2, 'always', 0],
        'subject-case': [
            2,
            'always',
            ['sentence-case', 'lower-case']
        ],
        'subject-empty': [2, 'never'],
        'subject-full-stop': [2, 'never', '.'],
        'signed-off-by': [2, 'never'],
        'type-case': [2, 'always', 'lower-case'],
        'type-empty': [2, 'never'],
        'type-enum': [
            2,
            'always',
            [
                'build',
                'chore',
                'ci',
                'docs',
                'feat',
                'fix',
                'perf',
                'refactor',
                'revert',
                'style',
                'test',
                'thing',
                'user',
                'amqp',
                'handler',
                'server',
                'config',
                'github',
            ]
        ]
    }
};
