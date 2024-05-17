import { defineStore } from "pinia";
import { Order } from "@/class/order";
import {socket} from "@/socket";


export const useOrderStore = defineStore("orders",{
    state: () => ({
        orders: [],
    }),

    actions: {
        bindEvents(){
            socket.onopen = () => {
                console.log("Successfully connected!");
            }
            socket.onmessage = (event) => {
                let json = JSON.parse(event.data)
                for (const [key, value] of Object.entries(json)){
                    switch (key){
                        case "Orders": {
                            const ordersArray = [...Order.generateOrder(value)]
                            console.log(ordersArray.sort((a,b) => a.id - b.id ))
                            this.orders = ordersArray.sort((a,b) => a.id - b.id )
                            break
                        }
                        case "Order": {
                            const order = Order.deserialize(value);
                            console.log(order);
                            const existingOrderIndex = this.orders.findIndex((temp) => {
                                return temp.id === order.id;
                            });
                            if (existingOrderIndex) {
                                this.orders[existingOrderIndex] = order
                            }
                            break
                        }
                        default:
                            console.log("Error during parse JSON")
                    }
                }
            }
        },
        addNewOrder(room, camera, startDate, endDate){
            let new_order = Order.createOrder(room, camera, startDate, endDate)
            if(socket.readyState === WebSocket.OPEN){
                socket.send(Order.serialize(new_order))
            }
        }
    }
})