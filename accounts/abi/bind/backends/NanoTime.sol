pragma solidity 0.5.11;

contract NanoTime { // a wrapper around the precompile.
    function timestamp() view external returns (uint256 result) {
        uint256 number = block.number;
        assembly {
            let m := mload(0x40)
            mstore(m, number)
            if iszero(staticcall(gas, 0x010001, m, 0x20, m, 0x20)) {
                revert(0, 0)
            }
            result := mload(m)
        }
    }
}
