package util

func Coalesce[T comparable](values ...T) T {
    var zero T
    for _, v := range values {
        if v != zero {
            return v
        }
    }
    return zero
}