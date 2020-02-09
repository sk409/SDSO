export function truncate(str, length) {
  if (str.length <= length) {
    return str;
  }
  return str.substr(0, length) + "...";
}

export function count(sequence, target) {
  return sequence.split(target).length - 1;
}

export function setupTest(test) {
  const failed = "failed";
  const failedColor = "rgb(220, 102, 97)";
  const running = "running";
  const runningColor = "rgb(130, 209, 226)";
  const success = "success";
  const successColor = "rgb(107, 197, 143)";
  test.status = test.results.find(result => result.status.text === failed)
    ? failed
    : test.results.find(result => result.status.text === running)
    ? running
    : success;
  test.color =
    test.status === failed
      ? failedColor
      : test.status === running
      ? runningColor
      : successColor;
  test.results.forEach(result => {
    result.color =
      result.status.text === failed
        ? failedColor
        : result.status.text === running
        ? runningColor
        : successColor;
  });
}

class DateFormatter {
  ago(str) {
    const date = new Date(str);
    const now = new Date();
    const seconds = Math.max(
      0,
      Math.floor((now.getTime() - date.getTime()) / 1000)
    );
    if (seconds === 0) {
      return "たった今";
    }
    if (seconds < 60) {
      return seconds + "秒前";
    }
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) {
      return minutes + "分前";
    }
    const hours = Math.floor(minutes / 60);
    if (hours < 24) {
      return hours + "時間前";
    }
    const d = Math.floor(hours / 24);
    if (d < 7) {
      return d + "日前";
    }
    const week = Math.floor(d / 7);
    if (week < 4) {
      return week + "週間前";
    }
    return this.default(str);
  }
  default(str) {
    const p = n => {
      const s = `${n}`;
      return s.length < 2 ? "0" + s : s;
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
