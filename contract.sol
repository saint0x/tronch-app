/* SPDX-License-Identifier: MIT */
pragma solidity ^0.8.0;

pragma solidity ^0.8.0;

contract EscrowService {
    address public client;
    string public clientEmail;
    uint public expirationDate;
    uint public paymentAmount = 1.50 ether;
    address public seller;
    bool public isEscrowInitiated;
    bool public isReceived;
    bool public isDisputed;
    
    event EscrowInitiated(address client, address seller, uint amount, uint timestamp);
    event ReceiptConfirmed(address receiver, uint amount, uint timestamp);
    event EscrowDisputed(address disputant, uint timestamp);
    event DisputeResolved(address resolver, uint timestamp);
    
    modifier onlyClient() {
        require(msg.sender == client, "Only the client can access this function");
        _;
    }
    
    modifier onlySeller() {
        require(msg.sender == seller, "Only the seller can access this function");
        _;
    }
    
    constructor(address _client, string memory _clientEmail) {
        client = _client;
        clientEmail = _clientEmail;
        expirationDate = block.timestamp + 30 days;
    }

    function initiateEscrow() external payable onlyClient {
        require(msg.value == paymentAmount, "Payment amount must be 1.50 ETH");
        isEscrowInitiated = true;
        seller = msg.sender;
        emit EscrowInitiated(client, seller, msg.value, block.timestamp);
    }
    
    function confirmReceipt() external onlyClient {
        require(isEscrowInitiated, "Escrow has not been initiated");
        require(!isDisputed, "Escrow is being disputed");
        seller.transfer(address(this).balance);
        isReceived = true;
        emit ReceiptConfirmed(seller, address(this).balance, block.timestamp);
    }
    
    function disputeEscrow() external onlyClient {
        require(isEscrowInitiated, "Escrow has not been initiated");
        require(!isDisputed, "Escrow is already being disputed");
        isDisputed = true;
        emit EscrowDisputed(client, block.timestamp);
    }
    
    function resolveDispute() external onlySeller {
        require(isDisputed, "No dispute to resolve");
        isDisputed = false;
        isReceived = false;
        client.transfer(address(this).balance);
        emit DisputeResolved(seller, block.timestamp);
    }
    
    function getTransactionDetails() external view returns (address, address, uint, bool, bool, bool) {
        return (client, seller, paymentAmount, isEscrowInitiated, isReceived, isDisputed);
    }
}