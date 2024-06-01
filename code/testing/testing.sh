# # test 1 node
# go run testing_binary.go 300 &
# wait
# mv log.txt log-1Node.txt

# test 4 nodes
# go run testing_binary.go 75 &
# go run testing_binary.go 75 &
# go run testing_binary.go 75 &
# go run testing_binary.go 75 &
# wait
# mv log.txt log-4Node.txt

# test 8 nodes
go run testing_binary.go 38 &
go run testing_binary.go 38 &
go run testing_binary.go 38 &
go run testing_binary.go 38 &
go run testing_binary.go 37 &
go run testing_binary.go 37 &
go run testing_binary.go 37 &
go run testing_binary.go 37 &
wait
# mv log.txt log-8Node.txt
