## Engine v1 Metering Sequence

The v1 engine does not make a standard deduction and so does not need to do any direct conversions to universal credits.

```mermaid
sequenceDiagram
    Engine->>Metering: Reserve [sets the universal balance]
    Note over Engine,Metering: standard deduction is in native units

    loop every step
        Engine->>Metering: GetMaxSpendForInvocation [universal]
        Metering-->>Engine: max universal

        rect rgb(0, 102, 255)
        Engine->>Metering: Deduct [universal]
        Note over Engine,Metering: earmark capability spends

        Engine->>Metering: CreditToSpendingLimits
        Metering-->>Engine: limits in native units
        Engine->>Capability: Execute [with native limits]
        Capability-->>Engine: response with native usage
        Engine->>Metering: Settle [native]
        end
    end
```


## Engine v2 Metering Sequence

The v2 engine makes a standard deduction for all workflow invocations where the native units are `Compute`. To use the 
same interface flow of `Deduct` and `Settle`, the native units must be converted to universal units. To accomplish this,
a special function must be available to convert native units to universal for deduct to function.

```mermaid
sequenceDiagram
    Engine->>Metering: Reserve [sets the universal balance]
    Note over Engine,Metering: standard deduction is in native units

    Engine->>Metering: ConvertToBalance [native]
    Metering-->>Engine: universal credits
    
    rect rgb(255, 102, 0)
    Engine->>Metering: Deduct [universal]
    Note over Engine,Metering: earmark standard compute
    
    loop every step
        Engine->>Metering: GetMaxSpendForInvocation [universal]
        Metering-->>Engine: max universal

        rect rgb(0, 102, 255)
        Engine->>Metering: Deduct [universal]
        Note over Engine,Metering: earmark capability spends

        Engine->>Metering: CreditToSpendingLimits
        Metering-->>Engine: limits in native units
        Engine->>Capability: Execute [with native limits]
        Capability-->>Engine: response with native usage
        Engine->>Metering: Settle [native]
        end
    end
    
    Engine->>Metering: Settle [native]
    end
```