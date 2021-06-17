package main

import (
	client2 "doccer/client"
)

func main() {
	client := client2.NewClient("http://localhost:8080")

	id1, _ := client.Register("Jacob", "abacaba")
	id2, _ := client.Register("Kurt", "qwerty")
	id3, _ := client.Register("Jordan", "zxc")

	println(id1, id2, id3)

	jwt1, _ := client.Login("Jacob", "abacaba")
	jwt2, _ := client.Login("Kurt", "qwerty")
	jwt3, _ := client.Login("Jordan", "zxc")

	groupId, _ := client.CreateGroup("Converge", jwt1)
	println(groupId)

	_ = client.AddMember(groupId, id2, jwt1)

	docId, _ := client.CreateDoc("Jane Doe", "Text", "none", jwt1)
	println(docId)

	_ = client.ChangeGroupAccess(docId, groupId, "absolute", jwt1)

	doc, _ := client.GetDoc(docId, jwt3)
	if doc != nil {
		println("This doc shouldn't be accessible")
	}

	_ = client.ChangeMemberAccess(docId, id3, "read", jwt2)

	doc, _ = client.GetDoc(docId, jwt3)
	println(doc.Text, doc.Lang, doc.LinterStatus)
}
