package exercise

import (
	"fmt"
	"slices"
)

type Item struct {
	Name string
	Type string
}

type Player struct {
	Name      string
	Inventory []Item
}

func (p *Player) PickUpItem(item Item) {
	p.Inventory = append(p.Inventory, item)
	fmt.Printf("Player %s picked up %s item.\n", p.Name, item.Name)
}

func (p *Player) DropItem(itemName string) {
	for i, item := range p.Inventory {
		if item.Name == itemName {
			p.Inventory = slices.Delete(p.Inventory, i, i+1)
			fmt.Printf("Player %s dropped %s item.\n", p.Name, itemName)
			return
		}
	}
	fmt.Printf("Player %s does not have %s item in inventory.\n", p.Name, itemName)
}

func (p *Player) UseItem(targetItem Item) {
	for i, item := range p.Inventory {
		if item.Name == targetItem.Name {
			switch {
			case item.Type == "potion":
				fmt.Printf("Player %s drank the %s potion.\n", p.Name, item.Name)
			case item.Type == "powerup":
				fmt.Printf("Player %s used the %s powerup.\n", p.Name, item.Name)
			default:
				fmt.Printf("Player %s cannot consume the %s item.\n", p.Name, item.Name)
			}

			p.Inventory = slices.Delete(p.Inventory, i, i+1) //. or
			// p.Inventory = append(p.Inventory[:i],p.Inventory[i+1:]...)
			return
		}
	}
	fmt.Printf("Player %s does not have %s in inventory.\n", p.Name, targetItem.Name)
}
