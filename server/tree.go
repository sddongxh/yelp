package main

import "fmt"

type Node struct { // Tree Node, this is for binary tree
	val   int
	left  *Node
	right *Node
}

func buildTree(nums []int) *Node {
	if len(nums) == 0 {
		return nil
	}
	k := len(nums) / 2
	root := &Node{val: nums[k], left: nil, right: nil}
	if k > 0 {
		root.left = buildTree(nums[0:k])
	}
	if k+1 < len(nums) {
		root.right = buildTree(nums[k+1:])
	}
	return root
}

func (r *Node) dfs(nums *[]int) {
	if r == nil {
		return
	}
	fmt.Println(r.val)
	*nums = append(*nums, r.val)
	r.left.dfs(nums)
	r.right.dfs(nums)
}

func dfs(r *Node, nums *[]int) {
	if r == nil {
		return
	}
	fmt.Println(r.val)
	*nums = append(*nums, r.val)
	dfs(r.left, nums)
	dfs(r.right, nums)
}

func main() {
	var nums []int = []int{1, 2, 3, 4, 5, 6, 7, 8}

	root := buildTree(nums)
	var all []int
	root.dfs(&all)
	fmt.Println(all)
	var all2 []int
	dfs(root, &all2)
	fmt.Println(all2)
}
