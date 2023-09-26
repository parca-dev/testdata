
#include <cstdio>
#include <pthread.h>

int __attribute__((noinline)) top1() {
  for (int i = 0; i < 1000; i++) {
  }

  return 0;
}

int __attribute__((noinline)) top2() {

  for (int i = 0; i < 1000; i++) {
  }

  return 0;
}

// ones
int __attribute__((noinline)) c1() { return top1(); }

int __attribute__((noinline)) b1() { return c1(); }

int __attribute__((noinline)) a1() { return b1(); }

// twos
int __attribute__((noinline)) c2() { return top2(); }

int __attribute__((noinline)) b2() { return c2(); }

int __attribute__((noinline)) a2() { return b2(); }

int __attribute__((noinline)) top1t() {
  while (true) {
    printf("1\n");
  }

  return 0;
}

int __attribute__((noinline)) top2t() {

  while (true) {
    printf("2\n");
  }

  return 0;
}

// ones
int __attribute__((noinline)) c1t() { return top1t(); }

int __attribute__((noinline)) b1t() { return c1t(); }

void *__attribute__((noinline)) a1t(void *) {
  b1t();
  return NULL;
}

// twos
int __attribute__((noinline)) c2t() { return top2t(); }

int __attribute__((noinline)) b2t() { return c2t(); }

void *__attribute__((noinline)) a2t(void *) {
  b2t();
  return NULL;
}

int main() {
  pthread_t thread_id;
  pthread_t thread_id2;

  pthread_create(&thread_id, NULL, a1t, NULL);
  pthread_create(&thread_id2, NULL, a2t, NULL);

  while (true) {
    a1();
    a2();
  }
  return 0;
}