package output

type Output interface {
    PrintResult()
    PrintTestResult()
}

type DefaultOutput struct {

}