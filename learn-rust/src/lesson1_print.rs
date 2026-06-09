fn format_print() {
    println!("占位: {}", "😄");
    println!("占位 1: {0}, 2: {1}, 1: {0}", 0, 1);
}

#[cfg(test)]
#[test]
fn test_format_print() {
    format_print();
}
