const isProd = process.env.NODE_ENV === 'production';

export const time = isProd ? () => { } : console.time
export const timeEnd = isProd ? () => { } : console.timeEnd
