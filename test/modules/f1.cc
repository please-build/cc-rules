module f1;
import q1;

int f(int n) {
  if (n == 0) {
    return 1;
  }
  return 1+q(n-1);
}
