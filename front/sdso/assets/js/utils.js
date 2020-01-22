class DateFormatter {
  default (str) {
    const p = (n) => {
      const s = `${n}`;
      return s.length < 2 ? "0" + s : s
    };
    const date = new Date(str);
    const y = date.getFullYear();
    const m = p(date.getMonth() + 1);
    const d = p(date.getDate());
    const h = p(date.getHours());
    const min = p(date.getMinutes());
    const s = p(date.getSeconds());
    return `${y}年${m}月${d}日 ${h}:${min}:${s}`;
  }
}
export const dateFormatter = new DateFormatter();


export function truncate(str, length) {
  if (str.length <= length) {
    return str;
  }
  return str.substr(0, length) + "...";
}
