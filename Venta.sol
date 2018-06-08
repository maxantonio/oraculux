pragma solidity ^0.4.18;

interface IMyContract {
    function balanceOf(address owner) external returns (uint256);

    function transfer(address to, uint256 amount) external returns (bool);
}

contract DXTSale {
    IMyContract public tokenContract;  // the token being sold
    uint256 public price;              // the price, in wei, per token
    address owner;
    uint256 public tokensSold;

    event Sold(address buyer, uint256 amount);

    function DXTSale(IMyContract _tokenContract) public {
        owner = msg.sender;
        tokenContract = _tokenContract;
        price = 0.1 finney;
    }

    function safeMultiply(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        } else {
            uint256 c = a * b;
            assert(c / a == b);
            return c;
        }
    }

    function() public payable {
        if (msg.value > 0) {
            uint numberOfTokens = msg.value / price;
            uint256 scaledAmount = safeMultiply(numberOfTokens,
                uint256(10) ** 2);
            require(tokenContract.balanceOf(this) >= scaledAmount);
            emit Sold(msg.sender, numberOfTokens);
            tokensSold += numberOfTokens;
            require(tokenContract.transfer(msg.sender, scaledAmount));
        }
    }

    function buyTokens() public payable {
        uint numberOfTokens = msg.value / price;
        uint256 scaledAmount = safeMultiply(numberOfTokens,
            uint256(10) ** 2);
        require(tokenContract.balanceOf(this) >= scaledAmount);
        emit Sold(msg.sender, numberOfTokens);
        tokensSold += numberOfTokens;
        require(tokenContract.transfer(msg.sender, scaledAmount));
    }

    function endSale() public {
        require(msg.sender == owner);
        // Send unsold tokens to the owner.
        require(tokenContract.transfer(owner, tokenContract.balanceOf(this)));
        msg.sender.transfer(address(this).balance);
    }
}