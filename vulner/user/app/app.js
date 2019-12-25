function sum(a, b) {
  for (let i = 0; i < 10000000000; ++i);
  return a + b + 2;
}
module.exports = sum;